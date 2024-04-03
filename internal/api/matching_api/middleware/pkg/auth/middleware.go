package auth

import "github.com/gin-gonic/gin"

func APIMiddleware(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Request.Header.Set("X-API-Key", apiKey)

		c.Next()
	}
}
