package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	migration "github.com/hervibest/be-yourmoments-backup/photo-svc/cmd/migrations"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/config"
	grpcHandler "github.com/hervibest/be-yourmoments-backup/photo-svc/internal/delivery/grpc"
	http "github.com/hervibest/be-yourmoments-backup/photo-svc/internal/delivery/http/controller"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/delivery/http/middleware"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/delivery/http/route"
	subscriber "github.com/hervibest/be-yourmoments-backup/photo-svc/internal/delivery/messaging"
	aiconsumer "github.com/hervibest/be-yourmoments-backup/photo-svc/internal/delivery/messaging/ai"
	uploadconsumer "github.com/hervibest/be-yourmoments-backup/photo-svc/internal/delivery/messaging/upload"
	producer "github.com/hervibest/be-yourmoments-backup/photo-svc/internal/gateway/messaging"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/consul"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/discovery"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/repository"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/usecase"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var logs = logger.New("main")
var (
	grpcServer *grpc.Server
	app        *fiber.App
)

func webServer(ctx context.Context) error {
	tp := config.InitTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	app = config.NewApp()
	app.Use(otelfiber.Middleware())
	tracer := otel.Tracer("fiber-server")

	serverConfig := config.NewServerConfig()
	dbConfig := config.NewPostgresDatabase()
	minioConfig := config.NewMinio()
	jetStreamConfig := config.NewJetStream()
	redisConfig := config.NewRedisClient()

	config.InitCreatorReviewStream(jetStreamConfig)
	config.InitUserStream(jetStreamConfig)
	config.DeleteAISimilarStream(jetStreamConfig, logs)
	config.InitAISimilarStream(jetStreamConfig)
	config.InitUploadPhotoStream(jetStreamConfig)

	registry, err := consul.NewRegistry(serverConfig.ConsulAddr, serverConfig.Name)
	if err != nil {
		logs.Error("Failed to create consul registry: " + err.Error())
		return err
	}

	GRPCserviceID := discovery.GenerateServiceID(serverConfig.Name + "-grpc")
	HTTPserviceID := discovery.GenerateServiceID(serverConfig.Name + "-http")

	grpcPortInt, _ := strconv.Atoi(serverConfig.GRPCPort)
	httpPortInt, _ := strconv.Atoi(serverConfig.HTTPPort)

	err = registry.RegisterService(ctx, serverConfig.Name+"-grpc", GRPCserviceID, serverConfig.GRPCInternalAddr, grpcPortInt, []string{"grpc"})
	if err != nil {
		logs.Error("Failed to register gRPC photo service to consul")
		return err
	}

	err = registry.RegisterService(ctx, serverConfig.Name+"-http", HTTPserviceID, serverConfig.HTTPInternalAddr, httpPortInt, []string{"http"})
	if err != nil {
		logs.Error("Failed to register HTTP photo service to consul")
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

	// aiAdapter, _ := adapter.NewAiAdapter(ctx, registry, logs)
	userAdapter, _ := adapter.NewUserAdapter(ctx, registry, logs)
	storageAdapter := adapter.NewStorageAdapter(minioConfig)
	CDNAdapter := adapter.NewCDNadapter()
	cacheAdapter := adapter.NewCacheAdapter(redisConfig)
	messaginAdapter := adapter.NewMessagingAdapter(jetStreamConfig)
	creatorProducer := producer.NewCreatorProducer(messaginAdapter, logs)
	photoProducer := producer.NewPhotoProducer(messaginAdapter, logs)

	customValidator := helper.NewCustomValidator()

	photoRepo, _ := repository.NewPhotoRepository(dbConfig)
	photoDetailRepo := repository.NewPhotoDetailRepository()
	facecamRepo := repository.NewFacecamRepository()
	userSimilarRepo := repository.NewUserSimilarRepository()
	exploreRepo, _ := repository.NewExploreRepository(dbConfig)
	creatorRepository, _ := repository.NewCreatorRepository(dbConfig)
	creatorDiscountRepository, _ := repository.NewCreatorDiscountRepository(dbConfig)
	bulkPhotoRepository := repository.NewBulkPhotoRepository()

	photoUseCase := usecase.NewPhotoUseCase(dbConfig, photoRepo, photoDetailRepo, userSimilarRepo, creatorRepository,
		bulkPhotoRepository, storageAdapter, CDNAdapter, logs)
	faceCamUseCase := usecase.NewFacecamUseCase(dbConfig, facecamRepo, userSimilarRepo, storageAdapter, logs)
	userSimilarPhotoUsecase := usecase.NewUserSimilarUsecase(dbConfig, photoRepo, photoDetailRepo, facecamRepo,
		userSimilarRepo, bulkPhotoRepository, userAdapter, photoProducer, logs)
	creatorUseCase := usecase.NewCreatorUseCase(dbConfig, creatorRepository, cacheAdapter, creatorProducer, logs)
	exploreUseCase := usecase.NewExploreUseCase(dbConfig, exploreRepo, photoRepo, CDNAdapter, tracer, logs)
	creatorDiscountUseCase := usecase.NewCreatorDiscountUseCase(dbConfig, creatorDiscountRepository, logs)
	checkoutUseCase := usecase.NewCheckoutUseCase(dbConfig, photoRepo, creatorRepository, creatorDiscountRepository, logs)

	userSimilarWorkerUC := usecase.NewUserSimilarWorkerUseCase(dbConfig, photoRepo, photoDetailRepo, facecamRepo,
		userSimilarRepo, bulkPhotoRepository, userAdapter, photoProducer, logs)

	photoUseCaseWorker := usecase.NewPhotoWorkerUseCase(dbConfig, photoRepo, photoDetailRepo, userSimilarRepo, creatorRepository,
		bulkPhotoRepository, storageAdapter, CDNAdapter, logs)

	facecamUseCaseWorker := usecase.NewFacecamUseCaseWorker(dbConfig, facecamRepo, userSimilarRepo, storageAdapter, logs)

	exploreController := http.NewExploreController(tracer, customValidator, exploreUseCase, logs)
	creatorDiscountController := http.NewCreatorDiscountController(creatorDiscountUseCase, customValidator, logs)
	healthCheckController := http.NewHealthCheckController()
	checkoutController := http.NewCheckoutController(checkoutUseCase, customValidator, logs)
	photoController := http.NewPhotoController(photoUseCase, customValidator, logs)
	authMiddleware := middleware.NewUserAuth(userAdapter, tracer, logs)
	creatorMiddleware := middleware.NewCreatorMiddleware(creatorUseCase, tracer, logs)

	aiSimilarConsumer := aiconsumer.NewAIConsumer(userSimilarWorkerUC, jetStreamConfig, logs)
	go func() {
		logs.Log("consume all ai simialr event beginning")
		if err := aiSimilarConsumer.ConsumeAllEvents(ctx); err != nil {
			logs.CustomError("failed to consume all event", err)
		}
	}()

	uploadConsumer := uploadconsumer.NewUploadConsumer(photoUseCaseWorker, facecamUseCaseWorker, jetStreamConfig, logs)
	go func() {
		logs.Log("consume all upload event beginning")
		if err := uploadConsumer.ConsumeAllEvents(ctx); err != nil {
			logs.CustomError("failed to consume all upload event", err)
		}
	}()

	creatorReviewSubscriber := subscriber.NewCreatorReviewSubscriber(jetStreamConfig, creatorUseCase, logs)
	go func() {
		if err := creatorReviewSubscriber.Start(ctx); err != nil {
			logs.Error(fmt.Sprintf("Subscriber error: %v", err))
		}
	}()

	userSubscriber := subscriber.NewUserSubscriber(jetStreamConfig, creatorUseCase, logs)
	go func() {
		if err := userSubscriber.Start(ctx); err != nil {
			logs.Error(fmt.Sprintf("Subscriber error: %v", err))
		}
	}()

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
		grpcHandler.NewPhotoGRPCHandler(grpcServer, photoUseCase, faceCamUseCase, userSimilarPhotoUsecase, creatorUseCase, checkoutUseCase)
		if err := grpcServer.Serve(l); err != nil {
			logs.Error(fmt.Sprintf("Failed to start gRPC server: %v", err))
		}
	}()

	routes := route.RouteConfig{
		App:                      app,
		ExploreController:        exploreController,
		HealthCheckController:    healthCheckController,
		CreatorDiscountControler: creatorDiscountController,
		PhotoController:          photoController,
		AuthMiddleware:           authMiddleware,
		CreatorMiddleware:        creatorMiddleware,
		CheckoutController:       checkoutController,
	}

	routes.Setup()
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

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if !config.IsLocal() {
		migration.Run()
	}

	if err := webServer(ctx); err != nil {
		logs.Error(err)
	}
}
