package database

import (
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis"
)

type RedisDataSources struct {
	RedisClient *redis.Client
}

// InitRC establishes connections to redis
func (rc *RedisDataSources) InitRC() error {
	log.Printf("Initializing redis connection \n")

	// Initialize redis connection
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	log.Printf("Connecting to Redis\n")
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "",
		DB:       0,
	})

	// verify redis connection

	_, err := rdb.Ping().Result()

	if err != nil {
		return fmt.Errorf("error connecting to redis: %w", err)
	}

	rc.RedisClient = rdb
	return nil
}
