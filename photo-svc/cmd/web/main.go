package main

import (
	"be-yourmoments/photo-svc/internal/adapter"
	"be-yourmoments/photo-svc/internal/config"
	grpcHandler "be-yourmoments/photo-svc/internal/delivery/grpc"
	http "be-yourmoments/photo-svc/internal/delivery/http/controller"
	"be-yourmoments/photo-svc/internal/delivery/http/middleware"
	"be-yourmoments/photo-svc/internal/delivery/http/route"
	"be-yourmoments/photo-svc/internal/helper"
	"log"
	"os"
	"os/signal"
	"syscall"

	"be-yourmoments/photo-svc/internal/helper/consul"
	"be-yourmoments/photo-svc/internal/helper/discovery"
	"be-yourmoments/photo-svc/internal/helper/logger"
	"be-yourmoments/photo-svc/internal/repository"
	"be-yourmoments/photo-svc/internal/usecase"
	"net"

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

func webServer() error {

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

	registry, err := consul.NewRegistry(serverConfig.ConsulAddr, serverConfig.Name)
	if err != nil {
		logs.Error("Failed to create consul registry for category service" + err.Error())
		return err
	}

	GRPCserviceID := discovery.GenerateServiceID(serverConfig.Name + "-grpc")
	HTTPserviceID := discovery.GenerateServiceID(serverConfig.Name + "-http")

	grpcPortInt, _ := strconv.Atoi(serverConfig.GRPCPort)
	httpPortInt, _ := strconv.Atoi(serverConfig.HTTPPort)

	ctx := context.Background()

	err = registry.RegisterService(ctx, serverConfig.Name+"-grpc", GRPCserviceID, serverConfig.GRPCAddr, grpcPortInt, []string{"grpc"})
	if err != nil {
		logs.Error("Failed to register gRPC photo service to consul")
		return err
	}

	err = registry.RegisterService(ctx, serverConfig.Name+"-http", HTTPserviceID, serverConfig.HTTPAddr, httpPortInt, []string{"http"})
	if err != nil {
		logs.Error("Failed to register category service to consuls")
		return err
	}

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
	defer registry.DeregisterService(ctx, GRPCserviceID)

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
	defer registry.DeregisterService(ctx, HTTPserviceID)

	aiAdapter, err := adapter.NewAiAdapter(ctx, registry)
	if err != nil {
		logs.Error(err)
	}

	userAdapter, err := adapter.NewUserAdapter(ctx, registry)
	if err != nil {
		logs.Error(err)
	}

	transactionAdapter, err := adapter.NewTransactionAdapter(ctx, registry)
	if err != nil {
		logs.Error(err)
	}

	logs.Log(fmt.Sprintf("Succsess connected http service at port: %v", serverConfig.HTTP))

	uploadAdapter := adapter.NewUploadAdapter(minioConfig)
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

	photoUsecase := usecase.NewPhotoUsecase(dbConfig, photoRepo, photoDetailRepo, userSimilarRepo, creatorRepository, aiAdapter, uploadAdapter, logs)
	faceCamUseCase := usecase.NewFacecamUseCase(dbConfig, facecamRepo, userSimilarRepo, aiAdapter, uploadAdapter, logs)
	userSimilarPhotoUsecase := usecase.NewUserSimilarUsecase(dbConfig, photoRepo, photoDetailRepo, facecamRepo, userSimilarRepo, logs)
	creatorUseCase := usecase.NewCreatorUseCase(dbConfig, creatorRepository, transactionAdapter, logs)
	exploreUseCase := usecase.NewExploreUseCase(dbConfig, exploreRepo, photoRepo, tracer, logs)
	creatorDiscountUseCase := usecase.NewCreatorDiscountUseCase(dbConfig, creatorDiscountRepository, logs)
	checkoutUseCase := usecase.NewCheckoutUseCase(dbConfig, photoRepo, creatorRepository, creatorDiscountRepository, logs)

	exploreController := http.NewExploreController(tracer, customValidator, exploreUseCase, logs)
	creatorDiscountController := http.NewCreatorDiscountController(creatorDiscountUseCase, customValidator, logs)
	healthCheckController := http.NewHealthCheckController()
	checkoutController := http.NewCheckoutController(checkoutUseCase, customValidator, logs)

	authMiddleware := middleware.NewUserAuth(userAdapter, tracer, logs)
	creatorMiddleware := middleware.NewCreatorMiddleware(creatorUseCase, tracer, logs)

	go func() {
		grpcServer := grpc.NewServer()
		reflection.Register(grpcServer)

		l, err := net.Listen("tcp", serverConfig.GRPC)
		if err != nil {
			logs.Error(fmt.Sprintf("Failed to listen: %v", err))
		}
		logs.Log(fmt.Sprintf("gRPC server started on %s", serverConfig.GRPC))
		defer l.Close()

		grpcHandler.NewPhotoGRPCHandler(grpcServer, photoUsecase, faceCamUseCase, userSimilarPhotoUsecase, creatorUseCase, checkoutUseCase)

		if err := grpcServer.Serve(l); err != nil {
			logs.Error(fmt.Sprintf("Failed to start gRPC category server: %v", err))
		}
	}()

	routeConfig := route.RouteConfig{
		App:                      app,
		ExploreController:        exploreController,
		HealthCheckController:    healthCheckController,
		CreatorDiscountControler: creatorDiscountController,
		AuthMiddleware:           authMiddleware,
		CreatorMiddleware:        creatorMiddleware,
		CheckoutController:       checkoutController,
	}
	routeConfig.Setup()

	logs.Log(fmt.Sprintf("Succsess connected http service at port: %v", serverConfig.HTTP))

	err = app.Listen(serverConfig.HTTP)

	if err != nil {
		logs.Error(fmt.Sprintf("Failed to start HTTP category server: %v", err))
		return err
	}
	return nil
}

func main() {
	if err := webServer(); err != nil {
		logs.Error(err)
	}

	logs.Log("Api gateway server started")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigchan
	logs.Log(fmt.Sprintf("Received signal: %s. Shutting down gracefully...", sig))
}
