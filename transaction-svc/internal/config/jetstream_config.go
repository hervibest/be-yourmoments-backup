package config

import (
	"fmt"
	"log"

	"github.com/hervibest/be-yourmoments-backup/transaction-svc/internal/helper/utils"
	"github.com/nats-io/nats.go"
)

func NewJetStream() nats.JetStreamContext {
	port := utils.GetEnv("NATS_JETSREAM_PORT")
	nc, err := nats.Connect(fmt.Sprintf("nats://localhost:%s", port))
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
