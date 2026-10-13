package handlers

import (
	"errors"
	"net/http"

	"github.com/ErenKarakus1/Distributed-LLM-Gateway/provider-service/internal/config"
	"github.com/ErenKarakus1/Distributed-LLM-Gateway/provider-service/internal/models"
	"github.com/ErenKarakus1/Distributed-LLM-Gateway/provider-service/internal/services"
	"github.com/ErenKarakus1/Distributed-LLM-Gateway/provider-service/internal/validations"
	"github.com/gin-gonic/gin"
)

func ProvideChatHandler(cfg config.Config) gin.HandlerFunc {
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

		resp, err := services.ProvideChat(req, cfg)
		if err != nil {
			if errors.Is(err, services.ErrUnsupportedProvider) {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if errors.Is(err, services.ErrOpenAIAPIKeyMissing) {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "openai api key is required"})
				return
			}
			if errors.Is(err, services.ErrOpenAIRequestFailed) {
				ctx.JSON(http.StatusBadGateway, gin.H{"error": "openai request failed"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		ctx.JSON(http.StatusOK, resp)
	}
}
