package config

import (
	"fmt"
	"log"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/utils"
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

func InitCreatorReviewStream(js nats.JetStreamContext) {
	_, err := js.AddStream(&nats.StreamConfig{
		Name:     "CREATOR_REVIEW_STREAM",
		Subjects: []string{"creator.review.updated"},
		Storage:  nats.FileStorage,
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		log.Fatalf("failed to create stream: %v", err)
	}
}

func InitUserStream(js nats.JetStreamContext) {
	_, err := js.AddStream(&nats.StreamConfig{
		Name:     "USER_STREAM",
		Subjects: []string{"user.created"},
		Storage:  nats.FileStorage,
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		log.Fatalf("failed to create stream: %v", err)
	}
}

func DeleteAISimilarStream(js nats.JetStreamContext, log *logger.Log) {
	err := js.DeleteStream("AI_SIMILAR_STREAM")
	if err != nil {
		log.CustomError("failed to delete stream", err)
		return
	}
	log.Log("successfully deleted AI_SIMILAR_STREAM")
}

func InitAISimilarStream(js nats.JetStreamContext) {
	_, err := js.AddStream(&nats.StreamConfig{
		Name:     "AI_SIMILAR_STREAM",
		Subjects: []string{"ai.bulk.photo", "ai.single.facecam", "ai.single.photo"},
		Storage:  nats.FileStorage,
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		log.Fatalf("failed to create stream: %v", err)
	}
}

func InitUploadPhotoStream(js nats.JetStreamContext) {
	_, err := js.AddStream(&nats.StreamConfig{
		Name:     "UPLOAD_PHOTO_STREAM",
		Subjects: []string{"upload.bulk.photo", "upload.single.facecam", "upload.single.photo", "upload.update.photo"},
		Storage:  nats.FileStorage,
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		log.Fatalf("failed to create stream: %v", err)
	}
}
