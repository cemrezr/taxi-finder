package matching_api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	circuitbreaker "taxi-finder/internal/circuitbreaker"
	"taxi-finder/internal/database/mongodb"
	"taxi-finder/utils"
)

func SetupMatchingApiRoutes(r *gin.Engine, client *mongodb.Client, cb *circuitbreaker.CircuitBreaker) {

	r.GET("/match", MatchHandler(client, cb))

}

func MatchHandler(client *mongodb.Client, cb *circuitbreaker.CircuitBreaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the circuit is open
		if cb.AllowRequest() {
			// Query parameters: customer location
			customerLatitude := c.Query("latitude")
			customerLongitude := c.Query("longitude")
			if customerLatitude == "" || customerLongitude == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Missing customer location"})
				return
			}

			// Convert customer location to float64
			customerLat, err := utils.ParseCoordinate(customerLatitude)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
				return
			}
			customerLng, err := utils.ParseCoordinate(customerLongitude)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
				return
			}

			// Get all drivers from the database
			drivers, err := client.FindDrivers()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}

			// Find nearest driver among all drivers
			nearestDriver, err := utils.FindNearestDriver(drivers, customerLat, customerLng)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}

			// Return the nearest driver
			c.JSON(http.StatusOK, nearestDriver)
		} else {
			// If circuit is open, return error
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service Unavailable"})
		}
	}
}
