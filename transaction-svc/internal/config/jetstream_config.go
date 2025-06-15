package config

import (
	"fmt"
	"log"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/utils"
	"github.com/nats-io/nats.go"
)

func NewJetStream() nats.JetStreamContext {
	host := utils.GetEnv("NATS_HOST")
	port := utils.GetEnv("NATS_PORT")
	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%s", host, port))
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Failed to get JetStream context: %v", err)
	}

	log.Println("Successfully connected to nats jetstream")

	return js
}

func InitCreatorStream(js nats.JetStreamContext) {
	err := js.DeleteStream("Creator_STREAM")
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "CREATOR_STREAM",
		Subjects: []string{"creator.created"},
		Storage:  nats.FileStorage,
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		log.Fatalf("failed to create stream: %v", err)
	}
}
