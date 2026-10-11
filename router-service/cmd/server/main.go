package main

import (
	"log"

	"github.com/ErenKarakus1/Distributed-LLM-Gateway/router-service/internal/config"
	"github.com/ErenKarakus1/Distributed-LLM-Gateway/router-service/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	cfg := config.LoadConfig()

	router.POST("/route/chat", handlers.RouteChatHandler(cfg.ProviderURL))

	if err := router.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}
