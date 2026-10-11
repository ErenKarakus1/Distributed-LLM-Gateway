package main

import (
	"log"

	"github.com/ErenKarakus1/Distributed-LLM-Gateway/provider-service/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/providers/chat", handlers.ProvideChatHandler())

	if err := router.Run(":8082"); err != nil {
		log.Fatal(err)
	}
}
