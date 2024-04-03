// driver-api/middleware/pkg/auth/auth.go
package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func APIMiddleware(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if apiKey != c.GetHeader("X-API-Key") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Next()
	}
}
