package config

import (
	"context"
	"log"

	"github.com/hervibest/be-yourmoments-backup/notification-svc/internal/helper/utils"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {

	host := utils.GetEnv("REDIS_HOST")
	port := utils.GetEnv("REDIS_PORT")
	address := host + ":" + port

	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})

	ctx := context.TODO()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	log.Println("✅ Redis client connected successfully...")

	return rdb
}
