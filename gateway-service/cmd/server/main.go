package main

import (
	"log"

	"github.com/ErenKarakus1/Distributed-LLM-Gateway/gateway-service/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/v1/chat/completions", handlers.CreateChatCompletionHandler())

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
