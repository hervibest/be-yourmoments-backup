package integration

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/config"
	grpcHandler "github.com/hervibest/be-yourmoments-backup/user-svc/internal/delivery/grpc"
	httphandler "github.com/hervibest/be-yourmoments-backup/user-svc/internal/delivery/http/controller"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/gateway/producer"
	"github.com/jmoiron/sqlx"

	"log"
	"net"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/delivery/http/middleware"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/delivery/http/route"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper/discovery"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper/discovery/consul"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/repository"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/usecase"

	"context"
	"fmt"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var logs = logger.New("INTEGRATION-TEST-USER-SVC")
var (
	grpcServer *grpc.Server
	app        *fiber.App
	dbConfig   *sqlx.DB
)

func webServer(ctx context.Context) error {
	app = config.NewApp()
	serverConfig := config.NewServerConfig()
	dbConfig = config.NewDB()
	minioConfig := config.NewMinio()
	redisConfig := config.NewRedisClient()
	firebaseConfig := config.NewFirebaseConfig()
	jetStreamConfig := config.NewJetStream()

	registry, err := consul.NewRegistry(serverConfig.ConsulAddr, serverConfig.Name)
	if err != nil {
		logs.Error("Failed to create consul registry for service" + err.Error())
		return err
	}

	GRPCserviceID := discovery.GenerateServiceID(serverConfig.Name + "-grpc")
	HTTPserviceID := discovery.GenerateServiceID(serverConfig.Name + "-http")

	grpcPortInt, _ := strconv.Atoi(serverConfig.GRPCPort)
	httpPortInt, _ := strconv.Atoi(serverConfig.HTTPPort)

	err = registry.RegisterService(ctx, serverConfig.Name+"-grpc", GRPCserviceID, serverConfig.GRPCInternalAddr, grpcPortInt, []string{"grpc"})
	if err != nil {
		logs.Error("Failed to register user service to consul")
		return err
	}

	err = registry.RegisterService(ctx, serverConfig.Name+"-http", HTTPserviceID, serverConfig.HTTPInternalAddr, httpPortInt, []string{"http"})
	if err != nil {
		logs.Error("Failed to register user service to consuls")
		return err
	}

	go func() {
		<-ctx.Done()
		logs.Log("Context canceled. Deregistering services...")
		registry.DeregisterService(context.Background(), GRPCserviceID)
		registry.DeregisterService(context.Background(), HTTPserviceID)

		logs.Log("Shutting down servers...")
		if err := app.Shutdown(); err != nil {
			logs.Error(fmt.Sprintf("Error shutting down Fiber: %v", err))
		}
		if grpcServer != nil {
			grpcServer.GracefulStop()
		}
		logs.Log("Successfully shutdown...")
	}()

	go startHealthCheckLoop(ctx, registry, GRPCserviceID, serverConfig.Name+"-grpc")
	go startHealthCheckLoop(ctx, registry, HTTPserviceID, serverConfig.Name+"-http")

	cacheAdapter := adapter.NewCacheAdapter(redisConfig)
	emailAdapter := adapter.NewEmailAdapter()
	googleTokenAdapter := adapter.NewGoogleTokenAdapter()
	jwtAdapter := adapter.NewJWTAdapter()
	securityAdapter := adapter.NewSecurityAdapter()
	uploadAdapter := adapter.NewUploadAdapter(minioConfig, redisConfig)
	realtimeChatAdapter := adapter.NewRealtimeChatAdapter(ctx, firebaseConfig, logs)
	authClientAdapter := adapter.NewAuthClientAdapter(firebaseConfig)
	cloudMessagingAdapter := adapter.NewCloudMessagingAdapter(firebaseConfig)
	perspectiveAdapter := adapter.NewPerspectiveAdapter()
	// photoAdapter, _ := adapter.NewPhotoAdapter(ctx, registry, logs)
	// transactionAdapter, _ := adapter.NewTransactionAdapter(ctx, registry, logs)
	messagingAdapter := adapter.NewMessagingAdapter(jetStreamConfig)

	userProducer := producer.NewUserProducer(messagingAdapter, logs)
	databaseAdapter := repository.NewDatabaseAdapter(dbConfig)

	customValidator := helper.NewCustomValidator()

	userRepository, err := repository.NewUserRepository(dbConfig)
	if err != nil {
		log.Fatalf(err.Error())
	}
	userProfileRepository, err := repository.NewUserProfileRepository(dbConfig)
	if err != nil {
		log.Fatalf(err.Error())
	}
	emailVerificationRepository, err := repository.NewEmailVerificationRepository(dbConfig)
	if err != nil {
		log.Fatalf(err.Error())
	}
	resetPasswordRepository, err := repository.NewResetPasswordRepository(dbConfig)
	if err != nil {
		log.Fatalf(err.Error())
	}

	userImageRepository, err := repository.NewUserImageRepository(dbConfig)
	if err != nil {
		log.Fatalf(err.Error())
	}

	userDeviceRepository := repository.NewUserDeviceRepository()

	authUseCase := usecase.NewAuthUseCase(databaseAdapter, userRepository, userProfileRepository, emailVerificationRepository, resetPasswordRepository,
		userDeviceRepository, googleTokenAdapter, emailAdapter, jwtAdapter, securityAdapter, cacheAdapter, realtimeChatAdapter,
		userProducer, logs)
	userUseCase := usecase.NewUserUseCase(databaseAdapter, userRepository, userProfileRepository, userImageRepository, uploadAdapter, cacheAdapter, logs)
	chatUseCase := usecase.NewChatUseCase(realtimeChatAdapter, authClientAdapter, cloudMessagingAdapter, perspectiveAdapter, logs)
	notificationUseCase := usecase.NewNotificationUseCase(databaseAdapter, redisConfig, userDeviceRepository, cloudMessagingAdapter, logs)
	authController := httphandler.NewAuthController(authUseCase, customValidator, logs)
	userController := httphandler.NewUserController(userUseCase, customValidator, logs)
	chatController := httphandler.NewChatController(chatUseCase, customValidator, logs)
	healthController := httphandler.NewHealthController()

	authMiddleware := middleware.NewUserAuth(authUseCase, customValidator, logs)

	routeConfig := route.RouteConfig{
		App:              app,
		AuthController:   authController,
		UserController:   userController,
		AuthMiddleware:   authMiddleware,
		ChatController:   chatController,
		HealthController: healthController,
	}

	go func() {
		grpcServer = grpc.NewServer()
		reflection.Register(grpcServer)
		l, err := net.Listen("tcp", serverConfig.GRPC)
		if err != nil {
			logs.Error(fmt.Sprintf("Failed to listen: %v", err))
			return
		}
		logs.Log(fmt.Sprintf("gRPC server started on %s", serverConfig.GRPC))
		defer l.Close()

		grpcHandler.NewUserGRPCHandler(grpcServer, authUseCase, notificationUseCase)

		if err := grpcServer.Serve(l); err != nil {
			logs.Error(fmt.Sprintf("Failed to start gRPC server: %v", err))
		}
	}()

	routeConfig.Setup()
	serverErrors := make(chan error, 1)
	go func() {
		logs.Log(fmt.Sprintf("Starting HTTP server at %s", serverConfig.HTTP))
		serverErrors <- app.Listen(serverConfig.HTTP)
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-serverErrors:
		return err
	}
}

func startHealthCheckLoop(ctx context.Context, registry *consul.Registry, serviceID, serviceName string) {
	failureCount := 0
	const maxFailures = 5
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := registry.HealthCheck(serviceID, serviceName)
			if err != nil {
				logs.Error(fmt.Sprintf("Failed to perform health check for %s: %v", serviceName, err))
				failureCount++
				if failureCount >= maxFailures {
					logs.Error(fmt.Sprintf("Max health check failures reached for %s. Exiting loop.", serviceName))
					return
				}
			} else {
				failureCount = 0
			}
			time.Sleep(2 * time.Second)
		}
	}
}
