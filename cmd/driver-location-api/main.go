// driver-location-api/main.go
package main

import (
	"github.com/gin-gonic/gin"
	"taxi-finder/internal/api/driver-api"
	"taxi-finder/internal/api/driver-api/middleware/pkg/auth"
	cb "taxi-finder/internal/circuitbreaker"
	"taxi-finder/internal/database/mongodb"
	"time"
)

func main() {
	client, err := mongodb.NewClient("mongodb://mongodb:27017")
	if err != nil {
	}

	cb := cb.NewCircuitBreaker(3, 5*time.Minute)

	r := gin.Default()

	r.Use(auth.APIMiddleware())

	driver_api.SetupRoutes(r, client, cb)
	r.Run(":8080")
}
