package config

import (
	"fmt"
	"log"

	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/helper/utils"
	"github.com/nats-io/nats.go"
)

func NewJetStream() nats.JetStreamContext {
	host := utils.GetEnv("NATS_HOST")
	port := utils.GetEnv("NATS_PORT")
	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%s", host, port))
	if err != nil {
		log.Println("Nats dsn : ", host, port)
		log.Fatalf("Failed to connect to NATS: %v", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Failed to get JetStream context: %v", err)
	}

	log.Println("Successfully connected to nats jetstream")

	return js
}

func InitPhotoStream(js nats.JetStreamContext, log logger.Log) {
	_, err := js.AddStream(&nats.StreamConfig{
		Name:     "PHOTO_STREAM",
		Subjects: []string{"photo.bulk", "photo.single.facecam", "photo.single.photo"},
		Storage:  nats.FileStorage,
	})

	if err != nil {
		log.CustomError("failed to cerate photo_stream ", err)
	}
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		log.CustomError("failed to setup photo stream", err)
	}
	log.Log("successfully created PHOTO_STREAM")
}

func DeletePhotoStream(js nats.JetStreamContext, log logger.Log) {
	err := js.DeleteStream("PHOTO_STREAM")
	if err != nil {
		log.CustomError("failed to delete stream", err)
		return
	}
	log.Log("successfully deleted PHOTO_STREAM")
}

func InitUserDeviceStream(js nats.JetStreamContext, log logger.Log) {
	_, err := js.AddStream(&nats.StreamConfig{
		Name:     "USER_DEVICE_STREAM",
		Subjects: []string{"user.device.created"},
		Storage:  nats.FileStorage,
	})

	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		log.CustomError("failed to setup photo stream", err)
	}
}
