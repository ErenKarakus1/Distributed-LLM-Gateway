package middleware

import (
	"net/http"

	"github.com/ErenKarakus1/Distributed-LLM-Gateway/gateway-service/internal/ratelimit"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimitMiddleware(rateLimitRequests int, rateLimitSeconds int, client *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKeyValue, exists := ctx.Get("api_key")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		apiKey, ok := apiKeyValue.(string)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		key := "rate_limit:" + apiKey
		allowed, err := ratelimit.Allow(ctx.Request.Context(), client, rateLimitRequests, rateLimitSeconds, key)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		if !allowed {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		ctx.Next()
	}
}
