package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func LoadEnvVariables() {
	err := godotenv.Load("/app/.env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	apiKey := os.Getenv("API_KEY")
	log.Println("API_KEY:", apiKey)
}

func APIMiddleware() gin.HandlerFunc {
	LoadEnvVariables()
	apiKey := os.Getenv("API_KEY")

	return func(c *gin.Context) {
		requestAPIKey := c.GetHeader("X-API-Key")

		if requestAPIKey != apiKey || requestAPIKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}
