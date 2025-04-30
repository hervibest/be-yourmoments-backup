package usecase

import (
	"context"
	"log"
	"testing"

	"firebase.google.com/go/messaging"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/adapter"
	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/config"
)

func TestSendFCM(t *testing.T) {
	firebaseApp := config.NewFirebaseConfig()
	// dbConfig := config.NewDB()
	// redisConfig := config.NewRedisClient()
	// deviceRepository := repository.NewUserDeviceRepository()

	cloudMessagingAdapter := adapter.NewCloudMessagingAdapter(firebaseApp)

	// notificationUseCase := usecase.NewNotificationUseCase(dbConfig, redisConfig, deviceRepository, cloudeMessagingAdapter, logger.New("TEST"))

	fcmToken := ""
	message := "Test ini adalah push notification"
	msg := &messaging.Message{
		Token: fcmToken,
		Notification: &messaging.Notification{
			Title: "Foto Mirip Terdeteksi",
			Body:  message, // Contoh: "Terdapat 5 foto yang mirip dengan Anda!"
		},
		Data: map[string]string{
			"type":    "similar_photo",
			"message": "message",
		},
	}

	_, err := cloudMessagingAdapter.Send(context.TODO(), msg)
	if err != nil {
		log.Printf("Error sending FCM to %s: %v", fcmToken, err)
	}
}

// type notificationUseCase struct {
// 	db                    *sqlx.DB
// 	redisClient           *redis.Client
// 	userDeviceRepository  repository.UserDeviceRepository
// 	cloudMessagingAdapter adapter.CloudMessagingAdapter

// 	logs *logger.Log
// }
