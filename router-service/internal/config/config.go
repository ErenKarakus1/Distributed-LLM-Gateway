package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ProviderURL string
}

func LoadConfig() Config {
	godotenv.Load()
	return Config{
		ProviderURL: getEnv("PROVIDER_URL", "http://localhost:8082/providers/chat"),
	}
}

func getEnv(s string, localDefault string) string {
	value := os.Getenv(s)
	if value == "" {
		return localDefault
	}
	return value
}
