package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	RouterURL         string
	GatewayAPIKey     string
	RedisURL          string
	RateLimitRequests int
	RateLimitSeconds  int
}

func LoadConfig() Config {
	godotenv.Load()
	rateLimitRequests, err := strconv.Atoi(getEnv("RATE_LIMIT_REQUESTS", "5"))
	if err != nil {
		rateLimitRequests = 5
	}
	rateLimitSeconds, err := strconv.Atoi(getEnv("RATE_LIMIT_SECONDS", "60"))
	if err != nil {
		rateLimitSeconds = 60
	}
	return Config{
		RouterURL:         getEnv("ROUTER_URL", "http://localhost:8081/route/chat"),
		GatewayAPIKey:     getEnv("GATEWAY_API_KEY", "dev-api-key"),
		RedisURL:          getEnv("REDIS_URL", "redis://localhost:6379"),
		RateLimitRequests: rateLimitRequests,
		RateLimitSeconds:  rateLimitSeconds,
	}
}

func getEnv(s string, localDefault string) string {
	value := os.Getenv(s)
	if value == "" {
		return localDefault
	}
	return value
}
