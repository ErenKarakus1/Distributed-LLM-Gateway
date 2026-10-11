package handlers

import (
	"errors"
	"net/http"

	"github.com/ErenKarakus1/Distributed-LLM-Gateway/provider-service/internal/models"
	"github.com/ErenKarakus1/Distributed-LLM-Gateway/provider-service/internal/services"
	"github.com/ErenKarakus1/Distributed-LLM-Gateway/provider-service/internal/validations"
	"github.com/gin-gonic/gin"
)

func ProvideChatHandler() gin.HandlerFunc {
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

		resp, err := services.ProvideChat(req)
		if err != nil {
			if errors.Is(err, services.ErrUnsupportedProvider) {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		ctx.JSON(http.StatusOK, resp)
	}
}
