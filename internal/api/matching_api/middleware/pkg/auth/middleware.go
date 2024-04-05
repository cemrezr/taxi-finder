// matching-api/auth/auth.go
package auth

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	err := godotenv.Load("/app/.env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
}

func AddAPIKeyMiddleware() gin.HandlerFunc {
	LoadEnvVariables()
	apiKey := os.Getenv("API_KEY")

	return func(c *gin.Context) {
		c.Set("API_KEY", apiKey)
		c.Next()
	}
}

func GetJWTSecret() []byte {
	LoadEnvVariables()
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	return jwtSecret
}

func JWTMiddleware() gin.HandlerFunc {
	jwtSecret := GetJWTSecret()

	return func(c *gin.Context) {
		tokenString := strings.Replace(c.GetHeader("Authorization"), "Bearer ", "", 1)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse token claims"})
			c.Abort()
			return
		}

		authenticated, exists := claims["authenticated"].(bool)
		if !exists || !authenticated {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not authenticated"})
			c.Abort()
			return
		}

		c.Next()
	}
}
