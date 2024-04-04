package models

type MatchResponse struct {
	DriverID string  `json:"driver_id"`
	Distance float64 `json:"distance"`
}
