package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthenticateMatchAPI(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestAPIKey := c.GetHeader("X-API-Key")

		if requestAPIKey != apiKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Next()
	}
}
