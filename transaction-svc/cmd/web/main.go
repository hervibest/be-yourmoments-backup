package main

import (
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/cmd/migration"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/config"
	grpcHandler "github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/delivery/grpc"
	http "github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/delivery/http/controller"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/delivery/messaging"
	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/delivery/scheduler"
	producer "github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/gateway/messaging"

	"log"
	"net"
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
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	var tracer = otel.Tracer("fiber-server")

	serverConfig := config.NewServerConfig()
	dbConfig := config.NewPostgresDatabase()
	midtransConfig := config.NewMidtransClient()
	redisConfig := config.NewRedisClient()
	goCronConfig := config.NewGocron()
	jetStreamConfig := config.NewJetStream()

	config.InitCreatorStream(jetStreamConfig)

	registry, err := consul.NewRegistry(serverConfig.ConsulAddr, serverConfig.Name)
	if err != nil {
		logs.Error("Failed to create consul registry for category service")
		return err
	}

	GRPCserviceID := discovery.GenerateServiceID(serverConfig.Name + "-grpc")
	HTTPserviceID := discovery.GenerateServiceID(serverConfig.Name + "-http")

	grpcPortInt, _ := strconv.Atoi(serverConfig.GRPCPort)
	httpPortInt, _ := strconv.Atoi(serverConfig.HTTPPort)

	err = registry.RegisterService(ctx, serverConfig.Name+"-grpc", GRPCserviceID, serverConfig.GRPCInternalAddr, grpcPortInt, []string{"grpc"})
	if err != nil {
		logs.Error("Failed to register gRPC transaction service to consul")
		return err
	}

	err = registry.RegisterService(ctx, serverConfig.Name+"-http", HTTPserviceID, serverConfig.HTTPInternalAddr, httpPortInt, []string{"http"})
	if err != nil {
		logs.Error("Failed to register category service to consul")
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

	photoAdapter, err := adapter.NewPhotoAdapter(ctx, registry, logs)
	if err != nil {
		logs.Error(err)
	}

	userAdapter, err := adapter.NewUserAdapter(ctx, registry, logs)
	if err != nil {
		logs.Error(err)
	}
	cacheAdapter := adapter.NewCacheAdapter(redisConfig)
	paymentAdapter := adapter.NewPaymentAdapter(midtransConfig, cacheAdapter, logs)
	messagingAdapter := adapter.NewMessagingAdapter(jetStreamConfig)
	// jwtAdapter := adapter.NewJWTAdapter()

	transactionProducer := producer.NewTransactionProducer(cacheAdapter, messagingAdapter, logs)

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

	walletUseCase := usecase.NewWalletUseCase(walletRepository, cacheAdapter, dbConfig, logs)
	bankUseCase := usecase.NewBankUseCase(dbConfig, bankRepository, logs)
	bankWalletUseCase := usecase.NewBankWalletUseCase(dbConfig, bankWalletRepoistory, bankRepository, logs)
	reviewUseCase := usecase.NewReviewUseCase(transactionDetailRepo, creatorReviewRepo, transactionProducer, dbConfig, logs)
	withdrawalUseCase := usecase.NewWithdrawalUseCase(dbConfig, withdrawalRepository, walletRepository, logs)
	transactionWalletUC := usecase.NewTransactionWalletUseCase(dbConfig, transactionWalletRepo, logs)
	cancelationUseCase := usecase.NewCancelationUseCase(dbConfig, transactionRepo, logs)
	schedulerUseCase := usecase.NewSchedulerUseCase(dbConfig, transactionRepo, transactionUseCase, cancelationUseCase, paymentAdapter, logs)

	transactionController := http.NewTransactionController(transactionUseCase, customValidator, logs)
	bankController := http.NewBankController(bankUseCase, customValidator, logs)
	bankWalletController := http.NewBankWalletController(bankWalletUseCase, customValidator, logs)
	reviewController := http.NewReviewController(reviewUseCase, customValidator, logs)
	withdarawlController := http.NewWithdrawalController(withdrawalUseCase, customValidator, logs)
	walletController := http.NewWalletController(walletUseCase, logs)
	transactionWalletCtrl := http.NewTransactionWalletController(transactionWalletUC, customValidator, logs)

	// authMiddlewareV2 := middleware.NewAuthMiddleware(userAdapter, logs, tracer, jwtAdapter, cacheAdapter)

	// authMiddleware := middleware.NewUserAuth(userAdapter, tracer, logs)
	creatorMiddleware := middleware.NewCreatorMiddleware(photoAdapter, tracer, logs)
	walletMiddleware := middleware.NewWalletMiddleware(walletUseCase, tracer, logs)

	go func() {
		grpcServer = grpc.NewServer()
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

	//TODO NEW - implement ctx for cancl orchestration
	schedulerRunner := scheduler.NewSchedulerRunner(goCronConfig, schedulerUseCase, logs)
	go func() {
		schedulerRunner.Start()
	}()

	creatorSubscriber := messaging.NewCreatorSubscriber(jetStreamConfig, walletUseCase, logs)
	go func() {
		if err := creatorSubscriber.Start(ctx); err != nil {
			logs.Error(fmt.Sprintf("Subscriber error: %v", err))
		}
	}()

	serverErrors := make(chan error, 1)
	route := route.NewRoute(app, transactionController, bankController, bankWalletController, reviewController,
		withdarawlController, walletController, transactionWalletCtrl, middleware.NewUserAuth(userAdapter, tracer, logs), creatorMiddleware, walletMiddleware)

	route.SetupRoute()
	app.Use(cors.New())

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
