package main

import (
	"github.com/hervibest/be-yourmoments-backup/user-svc/cmd/migration"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/config"
	grpcHandler "github.com/hervibest/be-yourmoments-backup/user-svc/internal/delivery/grpc"
	http "github.com/hervibest/be-yourmoments-backup/user-svc/internal/delivery/http/controller"

	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

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

var logs = logger.New("USER-SVC")

func webServer(ctx context.Context) error {
	app := config.NewApp()

	serverConfig := config.NewServerConfig()
	dbConfig := config.NewDB()
	minioConfig := config.NewMinio()
	redisConfig := config.NewRedisClient()
	firebaseConfig := config.NewFirebaseConfig()

	registry, err := consul.NewRegistry(serverConfig.ConsulAddr, serverConfig.Name)
	if err != nil {
		logs.Error("Failed to create consul registry for category service" + err.Error())
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
	}()

	go func() {
		failureCount := 0
		const maxFailures = 5
		for {
			err := registry.HealthCheck(GRPCserviceID, serverConfig.Name+"-grpc")
			if err != nil {
				logs.Error(fmt.Sprintf("Failed to perform health check for gRPC service: %v", err))
				failureCount++
				if failureCount >= maxFailures {
					logs.Error("Max health check failures reached for gRPC service. Exiting health check loop.")
					break
				}
			} else {
				failureCount = 0
			}
			time.Sleep(time.Second * 2)
		}
	}()

	go func() {
		failureCount := 0
		const maxFailures = 5
		for {
			err := registry.HealthCheck(HTTPserviceID, serverConfig.Name)
			if err != nil {
				logs.Error(fmt.Sprintf("Failed to perform health check: %v", err))
				failureCount++
				if failureCount >= maxFailures {
					logs.Error("Max health check failures reached for HTTP service. Exiting health check loop.")
					break
				}
			} else {
				failureCount = 0
			}
			time.Sleep(time.Second * 2)
		}
	}()

	cacheAdapter := adapter.NewCacheAdapter(redisConfig)
	emailAdapter := adapter.NewEmailAdapter()
	googleTokenAdapter := adapter.NewGoogleTokenAdapter()
	jwtAdapter := adapter.NewJWTAdapter()
	securityAdapter := adapter.NewSecurityAdapter()
	uploadAdapter := adapter.NewUploadAdapter(minioConfig, redisConfig)
	firestoreAdapter := adapter.NewFirestoreClientAdapter(firebaseConfig)
	authClientAdapter := adapter.NewAuthClientAdapter(firebaseConfig)
	cloudMessagingAdapter := adapter.NewCloudMessagingAdapter(firebaseConfig)
	perspectiveAdapter := adapter.NewPerspectiveAdapter()
	customValidator := helper.NewCustomValidator()

	photoAdapter, err := adapter.NewPhotoAdapter(ctx, registry, logs)
	if err != nil {
		log.Println(err)
	}

	transactionAdapter, err := adapter.NewTransactionAdapter(ctx, registry, logs)
	if err != nil {
		log.Println(err)
	}

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

	authUseCase := usecase.NewAuthUseCase(dbConfig, userRepository, userProfileRepository, emailVerificationRepository, resetPasswordRepository,
		userDeviceRepository, googleTokenAdapter, emailAdapter, jwtAdapter, securityAdapter, cacheAdapter, firestoreAdapter, photoAdapter, transactionAdapter, logs)
	userUseCase := usecase.NewUserUseCase(dbConfig, userRepository, userProfileRepository, userImageRepository, uploadAdapter, cacheAdapter, logs)
	chatUseCase := usecase.NewChatUseCase(firestoreAdapter, authClientAdapter, cloudMessagingAdapter, perspectiveAdapter, logs)
	notificationUseCase := usecase.NewNotificationUseCase(dbConfig, redisConfig, userDeviceRepository, cloudMessagingAdapter, logs)
	authController := http.NewAuthController(authUseCase, customValidator, logs)
	userController := http.NewUserController(userUseCase, customValidator, logs)
	chatController := http.NewChatController(chatUseCase, customValidator, logs)

	authMiddleware := middleware.NewUserAuth(authUseCase, customValidator, logs)

	routeConfig := route.RouteConfig{
		App:            app,
		AuthController: authController,
		UserController: userController,
		AuthMiddleware: authMiddleware,
		ChatController: chatController,
	}

	go func() {
		grpcServer := grpc.NewServer()
		reflection.Register(grpcServer)

		l, err := net.Listen("tcp", serverConfig.GRPC)
		if err != nil {
			logs.Error(fmt.Sprintf("Failed to listen: %v", err))
		}
		logs.Log(fmt.Sprintf("gRPC server started on %s", serverConfig.GRPC))
		defer l.Close()

		grpcHandler.NewUserGRPCHandler(grpcServer, authUseCase, notificationUseCase)

		if err := grpcServer.Serve(l); err != nil {
			logs.Error(fmt.Sprintf("Failed to start gRPC category server: %v", err))
		}
	}()

	routeConfig.Setup()
	logs.Log(fmt.Sprintf("Successfully connected http service at port: %v", serverConfig.HTTP))

	err = app.Listen(serverConfig.HTTP)

	if err != nil {
		logs.Error(fmt.Sprintf("Failed to start HTTP category server: %v", err))
		return err
	}
	return nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if !config.IsLocal() {
		migration.Run()
	}

	if err := webServer(ctx); err != nil {
		logs.Error(err)
	}

	logs.Log("Api gateway server started")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigchan
	logs.Log(fmt.Sprintf("Received signal: %s. Shutting down gracefully...", sig))
}
