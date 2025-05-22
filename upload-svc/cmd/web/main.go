package main

import (
	"os/signal"
	"syscall"

	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/config"
	grpcHandler "github.com/hervibest/be-yourmoments-backup/upload-svc/internal/delivery/grpc"
	http "github.com/hervibest/be-yourmoments-backup/upload-svc/internal/delivery/http/controller"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/delivery/http/middleware"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/delivery/http/route"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper/discovery"

	"net"

	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper/discovery/consul"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/usecase"

	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var logs = logger.New("main")
var (
	grpcServer *grpc.Server
	app        *fiber.App
)

func webServer(ctx context.Context) error {
	app = config.NewApp()
	serverConfig := config.NewServerConfig()
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

	err = registry.RegisterService(ctx, serverConfig.Name+"-grpc", GRPCserviceID, serverConfig.GRPCInternalAddr, grpcPortInt, []string{"grpc"})
	if err != nil {
		logs.Error("Failed to register gRPC upload service to consul")
		return err
	}

	err = registry.RegisterService(ctx, serverConfig.Name+"-http", HTTPserviceID, serverConfig.HTTPInternalAddr, httpPortInt, []string{"http"})
	if err != nil {
		logs.Error("Failed to register upload service to consuls")
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

	aiAdapter, err := adapter.NewAiAdapter(ctx, registry, logs)
	if err != nil {
		logs.Error(err)
	}

	photoAdapter, err := adapter.NewPhotoAdapter(ctx, registry, logs)
	if err != nil {
		logs.Error(err)
	}

	userAdapter, err := adapter.NewUserAdapter(ctx, registry, logs)
	if err != nil {
		logs.Error(err)
	}

	logs.Log(fmt.Sprintf("Success connected http service at port: %v", serverConfig.HTTP))

	storageAdapter := adapter.NewStorageAdapter(minioConfig)
	compressAdapter := adapter.NewCompressAdapter()
	customValidator := helper.NewCustomValidator()

	photoUsecase := usecase.NewPhotoUsecase(aiAdapter, photoAdapter, storageAdapter, compressAdapter, logs)
	photoController := http.NewPhotoController(photoUsecase, logs, customValidator)

	facecamUsecase := usecase.NewFacecamUseCase(aiAdapter, photoAdapter, storageAdapter, compressAdapter, logs)
	facecamController := http.NewFacecamController(facecamUsecase, logs)

	go func() {
		// gRPC server + reflection
		grpcServer = grpc.NewServer()
		reflection.Register(grpcServer)

		l, err := net.Listen("tcp", serverConfig.GRPC)
		if err != nil {
			logs.Error(fmt.Sprintf("Failed to listen: %v", err))
		}
		logs.Log(fmt.Sprintf("gRPC server started on %s", serverConfig.GRPC))
		defer l.Close()

		grpcHandler.NewPhotoGRPCHandler(grpcServer, photoUsecase)

		if err := grpcServer.Serve(l); err != nil {
			logs.Error(fmt.Sprintf("Failed to start gRPC category server: %v", err))
		}
	}()

	newUserMiddleware := middleware.NewUserAuth(userAdapter, logs)

	app.Use(cors.New(
		cors.ConfigDefault,
	))

	routeConfig := route.NewRouteConfig(app, photoController, facecamController, newUserMiddleware)
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

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := webServer(ctx); err != nil {
		logs.Error(err)
	}

	logs.Log("Received shutdown signal. Waiting for graceful shutdown...")
	<-ctx.Done()
	logs.Log("Shutdown complete.")
}
