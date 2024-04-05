package driver_api

import (
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"taxi-finder/internal/api/driver-api/database/mongodb"
	"taxi-finder/internal/circuitbreaker"
	"taxi-finder/internal/models"
	"taxi-finder/utils"
)

func SetupRoutes(r *gin.Engine, client *mongodb.Client, cb *circuitbreaker.CircuitBreaker) {

	r.GET("/drivers", GetDrivers(client))
	r.GET("/drivers/:id", GetDriverByID(client))
	r.POST("/drivers/create", CreateDriverLocation(client))
	r.POST("/match", MatchDriverAndRider(client, cb))
	r.POST("/upload-drivers", UploadDrivers(client))

}

func GetDrivers(client *mongodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		drivers, err := client.FindDrivers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(http.StatusOK, drivers)
	}
}

func GetDriverByID(client *mongodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		driverID := c.Param("id")

		driver, err := client.FindDriverByID(driverID)
		if err != nil {
			if err == mongodb.ErrDriverNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(http.StatusOK, driver)
	}
}

func UploadDrivers(client *mongodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		filePath := "cmd/driver-location-api/driver_coordinates/coordinates.csv"

		file, err := os.Open(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open CSV file"})
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		lines, err := reader.ReadAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CSV"})
			return
		}

		var drivers []*models.Driver

		for _, line := range lines {
			if len(line) != 2 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CSV format"})
				return
			}

			latitude, err := strconv.ParseFloat(line[0], 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
				return
			}

			longitude, err := strconv.ParseFloat(line[1], 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
				return
			}

			driver := &models.Driver{
				ID:        uuid.NewString(),
				Latitude:  latitude,
				Longitude: longitude,
			}
			drivers = append(drivers, driver)
		}

		err = client.InsertDriver(drivers)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save drivers to database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Drivers uploaded successfully"})
	}
}

func MatchDriverAndRider(client *mongodb.Client, cb *circuitbreaker.CircuitBreaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		var customerLocation models.GeoJSON
		if err := c.BindJSON(&customerLocation); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid GeoJSON format"})
			return
		}

		if !cb.AllowRequest() {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service Unavailable"})
			return
		}

		drivers, err := client.FindDrivers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		nearestDriver, distance, err := utils.FindNearestDriver(drivers, customerLocation.Coordinates[1], customerLocation.Coordinates[0])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"distance": distance,
			"driver": gin.H{
				"ID":        nearestDriver.ID,
				"Name":      nearestDriver.Name,
				"Latitude":  nearestDriver.Latitude,
				"Longitude": nearestDriver.Longitude,
			},
		})
	}
}

func CreateDriverLocation(client *mongodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Count int `json:"count"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if request.Count <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Count should be greater than zero"})
			return
		}

		var drivers []*models.Driver
		for i := 0; i < request.Count; i++ {
			longitude := rand.Float64()*180 - 90
			latitude := rand.Float64()*360 - 180

			driver := &models.Driver{
				ID:        uuid.NewString(),
				Latitude:  latitude,
				Longitude: longitude,
			}
			drivers = append(drivers, driver)
		}

		if err := client.InsertDriver(drivers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save drivers to database"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Driver locations created successfully", "count": len(drivers)})
	}
}
