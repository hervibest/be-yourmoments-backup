package main

import (
	"be-yourmoments/upload-svc/internal/adapter"
	"be-yourmoments/upload-svc/internal/config"
	grpcHandler "be-yourmoments/upload-svc/internal/delivery/grpc"
	"be-yourmoments/upload-svc/internal/delivery/http"
	"be-yourmoments/upload-svc/internal/delivery/http/middleware"
	"be-yourmoments/upload-svc/internal/helper"
	"be-yourmoments/upload-svc/internal/helper/discovery"
	"os"
	"os/signal"
	"syscall"

	"be-yourmoments/upload-svc/internal/helper/discovery/consul"
	"be-yourmoments/upload-svc/internal/helper/logger"
	"be-yourmoments/upload-svc/internal/usecase"
	"net"

	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/middleware/cors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var logs = logger.New("main")

func webServer() error {
	app := config.NewApp()
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

	ctx := context.Background()

	err = registry.RegisterService(ctx, serverConfig.Name+"-grpc", GRPCserviceID, serverConfig.GRPCAddr, grpcPortInt, []string{"grpc"})
	if err != nil {
		logs.Error("Failed to register gRPC book service to consul")
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

	aiAdapter, err := adapter.NewAiAdapter(ctx, registry, logs)
	if err != nil {
		logs.Error(err)
	}

	photoAdapter, err := adapter.NewPhotoAdapter(ctx, registry, logs)
	if err != nil {
		logs.Error(err)
	}

	userAdapter, err := adapter.NewUserAdapter(ctx, registry)
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
		grpcServer := grpc.NewServer()
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

	photoController.PhotoRoute(app, newUserMiddleware)
	facecamController.FacecamRoute(app, newUserMiddleware)
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
