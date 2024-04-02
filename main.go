package main

import (
	"context"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Driver struct {
	ID        string  `bson:"_id"`
	Name      string  `bson:"name"`
	Latitude  float64 `bson:"latitude"`
	Longitude float64 `bson:"longitude"`
}

type MatchResponse struct {
	DriverID string  `json:"driver_id"`
	Distance float64 `json:"distance"`
}

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongodb:27017"))
	if err != nil {
		log.Fatal("Error creating MongoDB client:", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	defer client.Disconnect(ctx)

	r := gin.Default()

	r.GET("/drivers", func(c *gin.Context) {
		collection := client.Database("test").Collection("drivers")

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			log.Printf("Error retrieving drivers: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		defer cursor.Close(ctx)

		var drivers []Driver
		if err := cursor.All(ctx, &drivers); err != nil {
			log.Printf("Error decoding drivers: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(http.StatusOK, drivers)
	})

	r.GET("/drivers/:id", func(c *gin.Context) {
		// Extract driver ID from URL parameter
		driverID := c.Param("id")

		// Access MongoDB collection
		collection := client.Database("test").Collection("drivers")

		// Set a longer context timeout for MongoDB operations
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		// Find the driver by ID
		var driver Driver
		err := collection.FindOne(ctx, bson.M{"_id": driverID}).Decode(&driver)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
				return
			}
			log.Printf("Error retrieving driver: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// Return the driver
		c.JSON(http.StatusOK, driver)
	})

	r.GET("/match", func(c *gin.Context) {
		customerLatitude := c.Query("latitude")
		customerLongitude := c.Query("longitude")
		if customerLatitude == "" || customerLongitude == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing customer location"})
			return
		}

		customerLat, err := parseCoordinate(customerLatitude)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
			return
		}
		customerLng, err := parseCoordinate(customerLongitude)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
			return
		}

		collection := client.Database("test").Collection("drivers")

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		// Find nearest driver using MongoDB's $near query
		cursor, err := collection.Find(ctx, bson.M{
			"latitude":  bson.M{"$exists": true},
			"longitude": bson.M{"$exists": true},
		})
		if err != nil {
			log.Printf("Error finding nearest driver: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		defer cursor.Close(ctx)

		var nearestDriver Driver
		var minDistance = math.MaxFloat64
		for cursor.Next(ctx) {
			var driver Driver
			if err := cursor.Decode(&driver); err != nil {
				log.Printf("Error decoding driver: %s", err)
				continue
			}
			distance := calculateDistance(customerLat, customerLng, driver.Latitude, driver.Longitude)
			if distance < minDistance {
				minDistance = distance
				nearestDriver = driver
			}
		}

		if nearestDriver.ID == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "No driver found"})
			return
		}

		response := MatchResponse{
			DriverID: nearestDriver.ID,
			Distance: minDistance,
		}
		c.JSON(http.StatusOK, response)
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error running server:", err)
	}
}

func parseCoordinate(coord string) (float64, error) {
	value, err := strconv.ParseFloat(coord, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371

	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadius * c

	return distance
}
