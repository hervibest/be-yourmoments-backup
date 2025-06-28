package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/config"
	subscriber "github.com/hervibest/be-yourmoments-backup/notification-svc/internal/delivery/messaging"
	consumer "github.com/hervibest/be-yourmoments-backup/notification-svc/internal/delivery/messaging/photo"
	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/repository"
	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/usecase"
)

var logs = logger.New("main")

func worker(ctx context.Context) error {
	dbConfig := config.NewPostgresDatabase()
	jetStreamConfig := config.NewJetStream()
	redisConfig := config.NewRedisClient()
	firebaseConfig := config.NewFirebaseConfig()

	config.DeletePhotoStream(jetStreamConfig, logs)
	config.InitPhotoStream(jetStreamConfig, logs)
	config.InitUserDeviceStream(jetStreamConfig, logs)

	go func() {
		<-ctx.Done()
		logs.Log("Context canceled. Deregistering services...")

		logs.Log("Shutting down servers...")

		logs.Log("Successfully shutdown...")
	}()

	cacheAdapter := adapter.NewCacheAdapter(redisConfig)
	databaseAdapter := repository.NewDatabaseAdapter(dbConfig)
	cloudMessagingAdapter := adapter.NewCloudMessagingAdapter(firebaseConfig)

	userDeviceRepo := repository.NewUserDeviceRepository()

	userDeviceUseCase := usecase.NewUserDeviceUseCase(databaseAdapter, userDeviceRepo, cacheAdapter, logs)
	notificationUseCase := usecase.NewNotificationUseCase(databaseAdapter, redisConfig, userDeviceRepo, cloudMessagingAdapter, logs)

	photoConsumer := consumer.NewPhotoConsumer(notificationUseCase, jetStreamConfig, logs)
	go func() {
		logs.Log("consume all photo event beginning")
		if err := photoConsumer.ConsumeAllEvents(ctx); err != nil {
			logs.CustomError("failed to consume all event", err)
		}
	}()

	userSubcriber := subscriber.NewUserSubscriber(jetStreamConfig, userDeviceUseCase, logs)
	go func() {
		if err := userSubcriber.Start(ctx); err != nil {
			logs.Error(fmt.Sprintf("Subscriber error: %v", err))
		}
	}()

	select {
	case <-ctx.Done():
		return nil
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := worker(ctx); err != nil {
		logs.Error(err)
	}
}
