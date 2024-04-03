package utils

import (
	"errors"
	"math"
	"strconv"
	"taxi-finder/internal/models"
)

func ParseCoordinate(coord string) (float64, error) {
	value, err := strconv.ParseFloat(coord, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Earth radius in kilometers
	const earthRadius = 6371

	// Convert latitude and longitude from degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Calculate differences
	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad

	// Haversine formula
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Calculate distance
	distance := earthRadius * c

	return distance
}

// FindNearestDriver fonksiyonunu oluştur
func FindNearestDriver(drivers []*models.Driver, customerLat, customerLng float64) (*models.Driver, error) {
	if len(drivers) == 0 {
		return nil, errors.New("no drivers available")
	}

	// Initialize variables to store nearest driver and minimum distance
	var nearestDriver *models.Driver
	var minDistance = math.MaxFloat64

	// Iterate through drivers to find the nearest one
	for _, driver := range drivers {
		distance := calculateDistance(customerLat, customerLng, driver.Latitude, driver.Longitude)
		if distance < minDistance {
			minDistance = distance
			nearestDriver = driver
		}
	}

	return nearestDriver, nil
}
