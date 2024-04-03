package driver_api

import (
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
	"strconv"
	circuitbreaker "taxi-finder/internal/circuitbreaker"
	"taxi-finder/internal/database/mongodb"
	"taxi-finder/internal/models"
)

func SetupRoutes(r *gin.Engine, client *mongodb.Client, cb *circuitbreaker.CircuitBreaker) {
	r.GET("/drivers", GetDriversHandler(client))
	r.GET("/drivers/:id", GetDriverByIDHandler(client))
	r.POST("/upload-drivers", UploadDriversHandler(client))

}

func GetDriversHandler(client *mongodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		drivers, err := client.FindDrivers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(http.StatusOK, drivers)
	}
}

func GetDriverByIDHandler(client *mongodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// URL parametresinden sürücü ID'sini al
		driverID := c.Param("id")

		// MongoDB'den sürücüyü ID'ye göre getir
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

func UploadDriversHandler(client *mongodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// CSV dosyasının yolu
		filePath := "cmd/driver-location-api/driver_coordinates/coordinates.csv"

		// CSV dosyasını aç
		file, err := os.Open(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open CSV file"})
			return
		}
		defer file.Close()

		// CSV dosyasını parse et
		reader := csv.NewReader(file)
		lines, err := reader.ReadAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CSV"})
			return
		}

		// Sürücüleri saklamak için bir dilim oluştur
		var drivers []*models.Driver

		// Parse edilen verileri Driver struct'ına dönüştür
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

			// Yeni bir Driver oluştur ve dilime ekle
			driver := &models.Driver{
				ID:        uuid.NewString(),
				Latitude:  latitude,
				Longitude: longitude,
			}
			drivers = append(drivers, driver)
		}

		// MongoDB'ye sürücüleri kaydet
		err = client.InsertDrivers(drivers)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save drivers to database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Drivers uploaded successfully"})
	}
}
