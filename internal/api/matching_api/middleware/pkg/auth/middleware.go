package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
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

func AddAPIKeyMiddleware() gin.HandlerFunc {
	LoadEnvVariables()
	apiKey := os.Getenv("API_KEY")

	return func(c *gin.Context) {
		c.Set("API_KEY", apiKey)
		c.Next()
	}
}
