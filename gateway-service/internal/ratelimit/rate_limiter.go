package ratelimit

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

func Allow(ctx context.Context, client *redis.Client, rateLimitRequests int, rateLimitSeconds int, key string) (bool, error) {
	value, err := client.Incr(ctx, key).Result()
	if err != nil {
		return false, errors.New("internal server error")
	}
	if value > int64(rateLimitRequests) {
		return false, nil
	}
	if value == 1 {
		duration := time.Duration(rateLimitSeconds) * time.Second
		ok, _ := client.Expire(ctx, key, duration).Result()
		if !ok {
			return false, errors.New("internal server error")
		}
	}
	return true, nil
}
