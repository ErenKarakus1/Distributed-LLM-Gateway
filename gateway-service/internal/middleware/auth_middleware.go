package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func APIKeyAuthMiddleware(expectedKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		expected := "Bearer " + expectedKey
		if authHeader != expected {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		apiKey := strings.TrimPrefix(authHeader, "Bearer ")
		ctx.Set("api_key", apiKey)
		ctx.Next()
	}
}
