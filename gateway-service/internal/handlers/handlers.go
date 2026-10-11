package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/ErenKarakus1/Distributed-LLM-Gateway/gateway-service/internal/models"
	"github.com/ErenKarakus1/Distributed-LLM-Gateway/gateway-service/internal/services"
	"github.com/ErenKarakus1/Distributed-LLM-Gateway/gateway-service/internal/validations"
	"github.com/gin-gonic/gin"
)

func CreateChatCompletionHandler(routerURL string) gin.HandlerFunc {
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

		start := time.Now()
		response, err := services.CreateChatCompletion(req, routerURL)
		duration := time.Since(start)
		if err != nil {
			var upstreamError models.UpstreamError
			if errors.As(err, &upstreamError) {
				log.Printf("chat completion request failed status=%v duration_ms=%v error=%v", upstreamError.StatusCode, duration.Milliseconds(), upstreamError.Error())
				ctx.JSON(upstreamError.StatusCode, gin.H{"error": upstreamError.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		log.Printf("completed chat completion request request_id=%v status=%v duration_ms=%v", response.RequestID, http.StatusOK, duration.Milliseconds())
		ctx.JSON(http.StatusOK, response)
	}
}
