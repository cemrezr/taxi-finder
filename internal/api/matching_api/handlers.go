package matching_api

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	circuitbreaker "taxi-finder/internal/circuitbreaker"
	"taxi-finder/internal/database/mongodb"
	"taxi-finder/internal/models"
	"taxi-finder/utils"
)

func SetupMatchingApiRoutes(r *gin.Engine, client *mongodb.Client, cb *circuitbreaker.CircuitBreaker) {

	r.GET("/nearest-driver", GetNearestDriver(client, cb))

}

func GetNearestDriver(client *mongodb.Client, cb *circuitbreaker.CircuitBreaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey, exists := c.Get("API_KEY")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "API_KEY not found in context"})
			return
		}
		_, ok := apiKey.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "API_KEY is not a string"})
			return
		}
		var customerLocation models.GeoJSON
		customerLatitude := c.Query("latitude")
		customerLongitude := c.Query("longitude")
		if customerLatitude == "" || customerLongitude == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing customer location"})
			return
		}
		lat, err := utils.ParseCoordinate(customerLatitude)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
			return
		}
		lng, err := utils.ParseCoordinate(customerLongitude)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
			return
		}

		customerLocation = models.GeoJSON{
			Type:        "Point",
			Coordinates: []float64{lng, lat},
		}

		requestBody := map[string]interface{}{
			"type":        customerLocation.Type,
			"coordinates": customerLocation.Coordinates,
		}

		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request body"})
			return
		}

		req, err := http.NewRequest("POST", "http://driver-location-api:8080/match", bytes.NewBuffer(jsonBody))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-Key", apiKey.(string))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get nearest driver from driver-api"})
			return
		}
		defer resp.Body.Close()

		c.Header("Content-Type", "application/json")
		c.Status(resp.StatusCode)
		_, err = io.Copy(c.Writer, resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy response body"})
			return
		}
	}
}
