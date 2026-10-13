package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenAIAPIKey  string
	OpenAIBaseURL string
	OpenAIModel   string
}

func LoadConfig() Config {
	godotenv.Load()
	return Config{
		OpenAIAPIKey:  getEnv("OPENAI_API_KEY", ""),
		OpenAIBaseURL: getEnv("OPENAI_BASE_URL", "https://api.openai.com/v1"),
		OpenAIModel:   getEnv("OPENAI_MODEL", "gpt-4.1-nano"),
	}
}

func getEnv(s string, localDefault string) string {
	value := os.Getenv(s)
	if value == "" {
		return localDefault
	}
	return value
}
