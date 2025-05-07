package main

import (
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/config"
	grpcHandler "github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/delivery/grpc"
	http "github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/delivery/http/controller"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/delivery/messaging"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/delivery/scheduler"
	producer "github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/gateway/messaging"

	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/delivery/http/middleware"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/delivery/http/route"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/repository"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/usecase"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/discovery"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/discovery/consul"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/logger"

	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	midtransConfig := config.NewMidtransClient()
	redisConfig := config.NewRedisClient()
	goCronConfig := config.NewGocron()

	registry, err := consul.NewRegistry(serverConfig.ConsulAddr, serverConfig.Name)
	if err != nil {
		logs.Error("Failed to create consul registry for category service")
		return err
	}

	GRPCserviceID := discovery.GenerateServiceID(serverConfig.Name + "-grpc")
	HTTPserviceID := discovery.GenerateServiceID(serverConfig.Name + "-http")

	grpcPortInt, _ := strconv.Atoi(serverConfig.GRPCPort)
	httpPortInt, _ := strconv.Atoi(serverConfig.HTTPPort)

	ctx := context.Background()

	err = registry.RegisterService(ctx, serverConfig.Name+"-grpc", GRPCserviceID, serverConfig.GRPCAddr, grpcPortInt, []string{"grpc"})
	if err != nil {
		logs.Error("Failed to register gRPC transaction service to consul")
		return err
	}

	err = registry.RegisterService(ctx, serverConfig.Name+"-http", HTTPserviceID, serverConfig.HTTPAddr, httpPortInt, []string{"http"})
	if err != nil {
		logs.Error("Failed to register category service to consul")
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

	photoAdapter, err := adapter.NewPhotoAdapter(ctx, registry)
	if err != nil {
		logs.Error(err)
	}

	userAdapter, err := adapter.NewUserAdapter(ctx, registry)
	if err != nil {
		logs.Error(err)
	}
	cacheAdapter := adapter.NewCacheAdapter(redisConfig)
	paymentAdapter := adapter.NewPaymentAdapter(midtransConfig, cacheAdapter, logs)
	transactionProducer := producer.NewTransactionProducer(cacheAdapter)

	customValidator := helper.NewCustomValidator()
	timeParserHelper := helper.NewTimeParserHelper(logs)

	transactionRepo := repository.NewTransactionRepository()
	transactionDetailRepo := repository.NewTransactionDetailRepository()
	transactionItemRepo := repository.NewTransactionItemRepository()

	bankRepository := repository.NewBankRepository()
	bankWalletRepoistory := repository.NewBankWalletRepository()
	creatorReviewRepo := repository.NewCreatorReviewRepository()
	withdrawalRepository := repository.NewWithdrawalRepository()
	walletRepository := repository.NewWalletRepository()
	transactionWalletRepo := repository.NewTransactionWalletRepository()

	transactionUseCase := usecase.NewTransactionUseCase(dbConfig, photoAdapter, transactionRepo,
		transactionItemRepo, transactionDetailRepo, walletRepository, transactionWalletRepo, paymentAdapter,
		cacheAdapter, timeParserHelper, transactionProducer, logs)

	walletUseCase := usecase.NewWalletUseCase(walletRepository, dbConfig, logs)
	bankUseCase := usecase.NewBankUseCase(dbConfig, bankRepository, logs)
	bankWalletUseCase := usecase.NewBankWalletUseCase(dbConfig, bankWalletRepoistory, logs)
	reviewUseCase := usecase.NewReviewUseCase(transactionDetailRepo, creatorReviewRepo, dbConfig, logs)
	withdrawalUseCase := usecase.NewWithdrawalUseCase(dbConfig, withdrawalRepository, walletRepository, logs)
	transactionWalletUC := usecase.NewTransactionWalletUseCase(dbConfig, transactionWalletRepo, logs)
	cancelationUseCase := usecase.NewCancelationUseCase(dbConfig, transactionRepo, logs)
	schedulerUseCase := usecase.NewSchedulerUseCase(dbConfig, transactionRepo, transactionUseCase, paymentAdapter, logs)

	transactionController := http.NewTransactionController(transactionUseCase, customValidator, logs)
	bankController := http.NewBankController(bankUseCase, customValidator, logs)
	bankWalletController := http.NewBankWalletController(bankWalletUseCase, customValidator, logs)
	reviewController := http.NewReviewController(reviewUseCase, customValidator, logs)
	withdarawlController := http.NewWithdrawalController(withdrawalUseCase, customValidator, logs)
	walletController := http.NewWalletController(walletUseCase, logs)
	transactionWalletCtrl := http.NewTransactionWalletController(transactionWalletUC, customValidator, logs)

	authMiddleware := middleware.NewUserAuth(userAdapter, tracer, logs)

	go func() {
		grpcServer := grpc.NewServer()
		reflection.Register(grpcServer)

		l, err := net.Listen("tcp", serverConfig.GRPC)
		if err != nil {
			logs.Error(fmt.Sprintf("Failed to listen: %v", err))
		}
		logs.Log(fmt.Sprintf("gRPC server started on %s", serverConfig.GRPC))
		defer l.Close()

		grpcHandler.NewTransactionGRPCHandler(grpcServer, walletUseCase)

		if err := grpcServer.Serve(l); err != nil {
			logs.Error(fmt.Sprintf("Failed to start gRPC category server: %v", err))
		}
	}()

	transactionSubscriber := messaging.NewTransactionSubsciber(cacheAdapter, cancelationUseCase, logs)
	go func() {
		transactionSubscriber.SubscribeTransactionExpire(ctx)
	}()

	schedulerRunner := scheduler.NewSchedulerRunner(goCronConfig, schedulerUseCase, logs)
	go func() {
		schedulerRunner.Start()
	}()

	route := route.NewRoute(app, transactionController, bankController, bankWalletController, reviewController,
		withdarawlController, walletController, transactionWalletCtrl, authMiddleware)

	route.SetupRoute()
	app.Use(cors.New())

	logs.Log(fmt.Sprintf("Successfully connected http service at port: %v", serverConfig.HTTP))

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

	logs.Log("Transaction service server started")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigchan
	logs.Log(fmt.Sprintf("Received signal: %s. Shutting down gracefully...", sig))
}
