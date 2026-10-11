package main

import (
	"log"

	"github.com/ErenKarakus1/Distributed-LLM-Gateway/gateway-service/internal/config"
	"github.com/ErenKarakus1/Distributed-LLM-Gateway/gateway-service/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	cfg := config.LoadConfig()

	router.POST("/v1/chat/completions", handlers.CreateChatCompletionHandler(cfg.RouterURL))

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
