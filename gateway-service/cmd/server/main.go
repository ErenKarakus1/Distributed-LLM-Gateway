package main

import (
	"log"

	"github.com/ErenKarakus1/Distributed-LLM-Gateway/gateway-service/internal/config"
	"github.com/ErenKarakus1/Distributed-LLM-Gateway/gateway-service/internal/handlers"
	"github.com/ErenKarakus1/Distributed-LLM-Gateway/gateway-service/internal/middleware"
	"github.com/ErenKarakus1/Distributed-LLM-Gateway/gateway-service/internal/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	cfg := config.LoadConfig()
	client := redis.NewClient(cfg.RedisURL)
	defer client.Close()

	router.POST("/v1/chat/completions", middleware.APIKeyAuthMiddleware(cfg.GatewayAPIKey), middleware.RateLimitMiddleware(cfg.RateLimitRequests, cfg.RateLimitSeconds, client), handlers.CreateChatCompletionHandler(cfg.RouterURL))

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
