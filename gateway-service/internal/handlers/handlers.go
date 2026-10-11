package handlers

import (
	"net/http"

	"github.com/ErenKarakus1/Distributed-LLM-Gateway/gateway-service/internal/models"
	"github.com/ErenKarakus1/Distributed-LLM-Gateway/gateway-service/internal/services"
	"github.com/ErenKarakus1/Distributed-LLM-Gateway/gateway-service/internal/validations"
	"github.com/gin-gonic/gin"
)

func CreateChatCompletionHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req models.CreateChatCompletionRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		req.Normalize()
		if err := validations.ValidateCreateChatCompletionRequest(req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, err := services.CreateChatCompletion(req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		ctx.JSON(http.StatusOK, response)
	}
}
