// matching-api/main.go
package main

import (
	"log"
	"taxi-finder/internal/api/matching_api"
	"taxi-finder/internal/api/matching_api/middleware"
	cb "taxi-finder/internal/circuitbreaker"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	circuitBreaker := cb.NewCircuitBreaker(3, 5*time.Minute)

	r := gin.Default()

	r.Use(middleware.AuthMiddleware())

	matching_api.SetupMatchingApiRoutes(r, circuitBreaker)

	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}
