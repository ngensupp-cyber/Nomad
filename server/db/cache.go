package db

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

func InitCache() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	RDB = redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	if err := RDB.Ping(Ctx).Err(); err != nil {
		log.Printf("Warning: Could not connect to Redis: %v. Real-time features might be limited.", err)
	} else {
		log.Println("Successfully connected to Redis")
	}
}
