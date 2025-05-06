package config

import (
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func NewGocron() gocron.Scheduler {
	s, err := gocron.NewScheduler(
		gocron.WithLocation(time.Local),
	)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	return s
}
