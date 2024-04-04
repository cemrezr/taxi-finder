// matching-api/main.go
package main

import (
	"github.com/gin-gonic/gin"
	"taxi-finder/internal/api/matching_api" // Adjusted import path
	"taxi-finder/internal/api/matching_api/middleware/pkg/auth"
	cb "taxi-finder/internal/circuitbreaker" // Adjusted import path
	"taxi-finder/internal/database/mongodb"  // Adjusted import path
	"time"
)

func main() {
	client, err := mongodb.NewClient("mongodb://mongodb:27017")
	if err != nil {
		panic(err)
	}

	cb := cb.NewCircuitBreaker(3, 5*time.Minute)

	r := gin.Default()

	r.Use(auth.AddAPIKeyMiddleware())

	matching_api.SetupMatchingApiRoutes(r, client, cb)

	r.Run(":8081")
}
