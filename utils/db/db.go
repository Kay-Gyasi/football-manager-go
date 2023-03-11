package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"time"
)

type IDatabase interface {
	Connect() (interface{}, error)
	Disconnect() error
	DoMigration() error
}

type Client struct {
	DB *mongo.Database

	coachCollection  *mongo.Collection
	playerCollection *mongo.Collection
	teamCollection   *mongo.Collection
}

func (c *Client) Connect() (interface{}, error) {
	var dbName string
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx)
	if err != nil {
		return nil, err
	}

	if os.Getenv("GIN_MODE") == "debug" {
		dbName = os.Getenv("DevDatabase")
	} else {
		dbName = os.Getenv("Database")
	}

	*c = Client{
		DB:               client.Database(dbName),
		coachCollection:  client.Database(dbName).Collection("coaches"),
		playerCollection: client.Database(dbName).Collection("players"),
		teamCollection:   client.Database(dbName).Collection("teams"),
	}

	return c.DB.Client(), nil
}

func (c *Client) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return c.DB.Client().Disconnect(ctx)
}
