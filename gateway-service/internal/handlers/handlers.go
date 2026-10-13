package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/ErenKarakus1/Distributed-LLM-Gateway/gateway-service/internal/models"
	"github.com/ErenKarakus1/Distributed-LLM-Gateway/gateway-service/internal/repository"
	"github.com/ErenKarakus1/Distributed-LLM-Gateway/gateway-service/internal/services"
	"github.com/ErenKarakus1/Distributed-LLM-Gateway/gateway-service/internal/validations"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateChatCompletionHandler(routerURL string, pool *pgxpool.Pool) gin.HandlerFunc {
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

		routerReq := models.RouterRequest{RequestID: uuid.New(), CreateChatCompletionRequest: req}

		start := time.Now()
		response, err := services.CreateChatCompletion(routerReq, routerURL)
		duration := time.Since(start)
		if err != nil {
			var upstreamError models.UpstreamError
			if errors.As(err, &upstreamError) {
				log.Printf("chat completion request failed status=%v duration_ms=%v error=%v", upstreamError.StatusCode, duration.Milliseconds(), upstreamError.Error())
				requestLog := services.CreateLogFail(routerReq, upstreamError.StatusCode, int(duration.Milliseconds()))
				if err := repository.InsertLog(ctx.Request.Context(), pool, requestLog); err != nil {
					log.Printf("Couldnt create log in db, request_id:%v", routerReq.RequestID)
				}
				ctx.JSON(upstreamError.StatusCode, gin.H{"error": upstreamError.Error()})
				return
			}
			requestLog := services.CreateLogFail(routerReq, 500, int(duration.Milliseconds()))
			if err := repository.InsertLog(ctx.Request.Context(), pool, requestLog); err != nil {
				log.Printf("Couldnt create log in db, request_id:%v", routerReq.RequestID)
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		responseLog := services.CreateLogSuccess(response, int(duration.Milliseconds()))
		if err := repository.InsertLog(ctx.Request.Context(), pool, responseLog); err != nil {
			log.Printf("Couldnt create log in db, request_id:%v", routerReq.RequestID)
		}
		log.Printf("completed chat completion request request_id=%v status=%v duration_ms=%v", response.RequestID, http.StatusOK, duration.Milliseconds())
		ctx.JSON(http.StatusOK, response)
	}
}
