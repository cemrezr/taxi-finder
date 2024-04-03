// driver-location-api/main.go
package main

import (
	"github.com/gin-gonic/gin"
	"taxi-finder/internal/api/driver-api"
	cb "taxi-finder/internal/circuitbreaker"
	"taxi-finder/internal/database/mongodb"
	"time"
)

func main() {
	// MongoDB ile bağlantı kur
	client, err := mongodb.NewClient("mongodb://mongodb:27017")
	if err != nil {
		// Hata kontrolü
	}

	// Circuit breaker'ı oluştur
	cb := cb.NewCircuitBreaker(3, 5*time.Minute)

	// HTTP sunucusunu başlat
	r := gin.Default()
	driver_api.SetupRoutes(r, client, cb) // Adjusted function call
	r.Run(":8080")
}
