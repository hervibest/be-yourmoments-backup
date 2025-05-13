package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	migration "github.com/hervibest/be-yourmoments-backup/photo-svc/cmd/migrations"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/config"
	grpcHandler "github.com/hervibest/be-yourmoments-backup/photo-svc/internal/delivery/grpc"
	http "github.com/hervibest/be-yourmoments-backup/photo-svc/internal/delivery/http/controller"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/delivery/http/middleware"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/delivery/http/route"
	subscriber "github.com/hervibest/be-yourmoments-backup/photo-svc/internal/delivery/messaging"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper"

	"net"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/consul"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/discovery"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/repository"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/usecase"

	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/contrib/otelfiber"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var logs = logger.New("main")

func webServer(ctx context.Context) error {

	tp := config.InitTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	app := config.NewApp()
	app.Use(otelfiber.Middleware())
	var tracer = otel.Tracer("fiber-server")

	serverConfig := config.NewServerConfig()
	dbConfig := config.NewPostgresDatabase()
	minioConfig := config.NewMinio()
	jetStreamConfig := config.NewJetStream()
	config.InitStream(jetStreamConfig)

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
		logs.Error("Failed to register gRPC photo service to consul")
		return err
	}

	err = registry.RegisterService(ctx, serverConfig.Name+"-http", HTTPserviceID, serverConfig.HTTPInternalAddr, httpPortInt, []string{"http"})
	if err != nil {
		logs.Error("Failed to register category service to consuls")
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

	aiAdapter, err := adapter.NewAiAdapter(ctx, registry, logs)
	if err != nil {
		logs.Error(err)
	}

	userAdapter, err := adapter.NewUserAdapter(ctx, registry, logs)
	if err != nil {
		logs.Error(err)
	}

	transactionAdapter, err := adapter.NewTransactionAdapter(ctx, registry, logs)
	if err != nil {
		logs.Error(err)
	}

	logs.Log(fmt.Sprintf("Successfully connected http service at port: %v", serverConfig.HTTP))

	storageAdapter := adapter.NewStorageAdapter(minioConfig)
	CDNAdapter := adapter.NewCDNadapter()
	customValidator := helper.NewCustomValidator()

	photoRepo, err := repository.NewPhotoRepository(dbConfig)
	if err != nil {
		log.Fatalln(err)
	}

	photoDetailRepo := repository.NewPhotoDetailRepository()
	facecamRepo := repository.NewFacecamRepository()
	userSimilarRepo := repository.NewUserSimilarRepository()
	exploreRepo, err := repository.NewExploreRepository(dbConfig)
	if err != nil {
		logs.Error(err)
	}

	creatorRepository, err := repository.NewCreatorRepository(dbConfig)
	if err != nil {
		log.Println(err)
	}

	creatorDiscountRepository, err := repository.NewCreatorDiscountRepository(dbConfig)
	if err != nil {
		log.Println(err)
	}
	bulkPhotoRepository := repository.NewBulkPhotoRepository()

	photoUseCase := usecase.NewPhotoUseCase(dbConfig, photoRepo, photoDetailRepo, userSimilarRepo, creatorRepository, bulkPhotoRepository, aiAdapter, storageAdapter, CDNAdapter, logs)
	faceCamUseCase := usecase.NewFacecamUseCase(dbConfig, facecamRepo, userSimilarRepo, aiAdapter, storageAdapter, logs)
	userSimilarPhotoUsecase := usecase.NewUserSimilarUsecase(dbConfig, photoRepo, photoDetailRepo, facecamRepo, userSimilarRepo, bulkPhotoRepository, userAdapter, logs)
	creatorUseCase := usecase.NewCreatorUseCase(dbConfig, creatorRepository, transactionAdapter, logs)
	exploreUseCase := usecase.NewExploreUseCase(dbConfig, exploreRepo, photoRepo, CDNAdapter, tracer, logs)
	creatorDiscountUseCase := usecase.NewCreatorDiscountUseCase(dbConfig, creatorDiscountRepository, logs)
	checkoutUseCase := usecase.NewCheckoutUseCase(dbConfig, photoRepo, creatorRepository, creatorDiscountRepository, logs)

	exploreController := http.NewExploreController(tracer, customValidator, exploreUseCase, logs)
	creatorDiscountController := http.NewCreatorDiscountController(creatorDiscountUseCase, customValidator, logs)
	healthCheckController := http.NewHealthCheckController()
	checkoutController := http.NewCheckoutController(checkoutUseCase, customValidator, logs)
	photoController := http.NewPhotoController(photoUseCase, customValidator, logs)

	authMiddleware := middleware.NewUserAuth(userAdapter, tracer, logs)
	creatorMiddleware := middleware.NewCreatorMiddleware(creatorUseCase, tracer, logs)

	creatorReviewSubscriber := subscriber.NewCreatorReviewSubscriber(jetStreamConfig, creatorUseCase, logs)
	go func() {
		creatorReviewSubscriber.Start(ctx)
	}()

	go func() {
		grpcServer := grpc.NewServer()
		reflection.Register(grpcServer)

		l, err := net.Listen("tcp", serverConfig.GRPC)
		if err != nil {
			logs.Error(fmt.Sprintf("Failed to listen: %v", err))
		}
		logs.Log(fmt.Sprintf("gRPC server started on %s", serverConfig.GRPC))
		defer l.Close()

		grpcHandler.NewPhotoGRPCHandler(grpcServer, photoUseCase, faceCamUseCase, userSimilarPhotoUsecase, creatorUseCase, checkoutUseCase)

		if err := grpcServer.Serve(l); err != nil {
			logs.Error(fmt.Sprintf("Failed to start gRPC category server: %v", err))
		}
	}()

	routeConfig := route.RouteConfig{
		App:                      app,
		ExploreController:        exploreController,
		HealthCheckController:    healthCheckController,
		CreatorDiscountControler: creatorDiscountController,
		PhotoController:          photoController,
		AuthMiddleware:           authMiddleware,
		CreatorMiddleware:        creatorMiddleware,
		CheckoutController:       checkoutController,
	}
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
