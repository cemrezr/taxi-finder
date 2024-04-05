// driver-location-api/main.go
package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"taxi-finder/internal/api/driver-api"
	"taxi-finder/internal/api/driver-api/database/mongodb"
	"taxi-finder/internal/api/driver-api/middleware"
	cb "taxi-finder/internal/circuitbreaker"
	"time"
)

func main() {
	client, err := mongodb.NewClient("mongodb://mongodb:27017")
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	cb := cb.NewCircuitBreaker(3, 5*time.Minute)

	r := gin.Default()

	r.Use(middleware.APIMiddleware())

	driver_api.SetupRoutes(r, client, cb)
	r.Run(":8080")
}
