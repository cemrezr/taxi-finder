// matching-api/main.go
package main

import (
	"github.com/gin-gonic/gin"
	"taxi-finder/internal/api/matching_api"  // Adjusted import path
	cb "taxi-finder/internal/circuitbreaker" // Adjusted import path
	"taxi-finder/internal/database/mongodb"  // Adjusted import path
	"time"
)

func main() {
	client, err := mongodb.NewClient("mongodb://mongodb:27017")
	if err != nil {
	}

	cb := cb.NewCircuitBreaker(3, 5*time.Minute)

	r := gin.Default()

	matchingAPIKey := "driver-api-key"

	//r.Use(auth.APIMiddleware(matchingAPIKey))

	matching_api.SetupMatchingApiRoutes(r, client, cb, matchingAPIKey)
	r.Run(":8081")
}
