package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RouterURL     string
	GatewayAPIKey string
}

func LoadConfig() Config {
	godotenv.Load()
	return Config{
		RouterURL:     getEnv("ROUTER_URL", "http://localhost:8081/route/chat"),
		GatewayAPIKey: getEnv("GATEWAY_API_KEY", "dev-api-key"),
	}
}

func getEnv(s string, localDefault string) string {
	value := os.Getenv(s)
	if value == "" {
		return localDefault
	}
	return value
}
