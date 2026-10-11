package middleware

import (
	"net/http"
	"time"

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
		value, err := client.Incr(ctx, key).Result()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		if value > int64(rateLimitRequests) {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		if value == 1 {
			duration := time.Duration(rateLimitSeconds) * time.Second
			ok, _ := client.Expire(ctx, key, duration).Result()
			if !ok {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
				return
			}
		}
		ctx.Next()
	}
}
