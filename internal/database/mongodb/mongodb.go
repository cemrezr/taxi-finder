package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"taxi-finder/internal/models"
)

var (
	ErrDriverNotFound = errors.New("driver not found")
)

type Client struct {
	db *mongo.Database
}

func NewClient(uri string) (*Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	db := client.Database("taxi-finder")

	return &Client{db: db}, nil
}

func (c *Client) FindDriverByID(id string) (*models.Driver, error) {
	var driver models.Driver
	err := c.db.Collection("drivers").FindOne(context.Background(), bson.M{"_id": id}).Decode(&driver)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrDriverNotFound
		}
		return nil, err
	}

	return &driver, nil
}

func (c *Client) FindDrivers() ([]*models.Driver, error) {
	cursor, err := c.db.Collection("drivers").Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var drivers []*models.Driver
	for cursor.Next(context.Background()) {
		var driver models.Driver
		err := cursor.Decode(&driver)
		if err != nil {
			return nil, err
		}
		drivers = append(drivers, &driver)
	}

	return drivers, nil
}

func (c *Client) InsertDrivers(drivers []*models.Driver) error {
	var documents []interface{}

	for _, driver := range drivers {
		documents = append(documents, driver)
	}

	_, err := c.db.Collection("drivers").InsertMany(context.Background(), documents)
	if err != nil {
		return err
	}

	return nil
}
