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
	// MongoDB ile bağlantı kur
	client, err := mongodb.NewClient("mongodb://mongodb:27017")
	if err != nil {
		// Hata kontrolü
	}

	// Circuit breaker'ı oluştur
	cb := cb.NewCircuitBreaker(3, 5*time.Minute)

	// HTTP sunucusunu başlat
	r := gin.Default()
	matching_api.SetupMatchingApiRoutes(r, client, cb) // Adjusted function call
	r.Run(":8081")
}
