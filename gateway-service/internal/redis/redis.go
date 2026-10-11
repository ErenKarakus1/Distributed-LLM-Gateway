package redis

import (
	"log"

	"github.com/redis/go-redis/v9"
)

func NewClient(redisURL string) *redis.Client {
	parsedURL, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatal(err)
	}
	return redis.NewClient(parsedURL)
}
