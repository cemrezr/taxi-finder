// matching-api/main.go
package main

import (
	"log"
	"taxi-finder/internal/api/matching_api"
	"taxi-finder/internal/api/matching_api/middleware/pkg/auth"
	cb "taxi-finder/internal/circuitbreaker"
	"taxi-finder/internal/database/mongodb"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	client, err := mongodb.NewClient("mongodb://mongodb:27017")
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	circuitBreaker := cb.NewCircuitBreaker(3, 5*time.Minute)

	r := gin.Default()

	r.Use(auth.AddAPIKeyMiddleware())
	r.Use(auth.JWTMiddleware())

	matching_api.SetupMatchingApiRoutes(r, client, circuitBreaker)

	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}
