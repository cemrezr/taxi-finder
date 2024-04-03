package models

type Driver struct {
	ID        string  `bson:"_id"`
	Name      string  `bson:"name"`
	Latitude  float64 `bson:"latitude"`
	Longitude float64 `bson:"longitude"`
}
