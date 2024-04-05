package utils

import (
	"errors"
	"taxi-finder/internal/models"
	"testing"
)

func TestParseCoordinate(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		err      error
	}{
		{"0", 0, nil},
		{"-90", -90, nil},
		{"180", 180, nil},
		{"abc", 0, errors.New("strconv.ParseFloat: parsing \"abc\": invalid syntax")},
	}

	for _, test := range tests {
		result, err := ParseCoordinate(test.input)
		if result != test.expected {
			t.Errorf("Expected: %f, got: %f", test.expected, result)
		}
		if (err == nil && test.err != nil) || (err != nil && test.err == nil) || (err != nil && test.err != nil && err.Error() != test.err.Error()) {
			t.Errorf("Expected error: %v, got: %v", test.err, err)
		}
	}
}

func TestFindNearestDriver(t *testing.T) {
	driver1 := &models.Driver{Latitude: 40.7128, Longitude: -74.0060}
	driver2 := &models.Driver{Latitude: 34.0522, Longitude: -118.2437}
	driver3 := &models.Driver{Latitude: 51.5074, Longitude: -0.1278}

	tests := []struct {
		drivers       []*models.Driver
		customerLat   float64
		customerLng   float64
		expected      *models.Driver
		expectedDist  float64
		expectedError error
	}{
		{[]*models.Driver{}, 40.7128, -74.0060, nil, 0, errors.New("no drivers available")},
		{[]*models.Driver{driver1, driver2, driver3}, 40.7128, -74.0060, driver1, 0, nil},
		{[]*models.Driver{driver1, driver2, driver3}, 34.0522, -118.2437, driver2, 0, nil},
		{[]*models.Driver{driver1, driver2, driver3}, 51.5074, -0.1278, driver3, 0, nil},
	}

	for _, test := range tests {
		driver, distance, err := FindNearestDriver(test.drivers, test.customerLat, test.customerLng)
		if driver != test.expected || distance != test.expectedDist || (err == nil && test.expectedError != nil) || (err != nil && test.expectedError == nil) || (err != nil && test.expectedError != nil && err.Error() != test.expectedError.Error()) {
			t.Errorf("For (%f, %f), expected: %v, %f, %v, got: %v, %f, %v", test.customerLat, test.customerLng, test.expected, test.expectedDist, test.expectedError, driver, distance, err)
		}
	}
}
