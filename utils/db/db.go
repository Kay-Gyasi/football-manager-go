package db

import (
	"context"
	"fmt"
	"football-manager-go/models"
	"football-manager-go/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

type IDatabase interface {
	Connect(ctx *context.Context) (interface{}, error)
	Disconnect(ctx *context.Context) error

	AddCoach(coach *models.Coach) (string, error)
	DeleteCoach(id primitive.ObjectID) error
	UpdateCoach(coach *models.Coach) error
	GetCoachById(id primitive.ObjectID) (models.Coach, error)
	GetCoachByFilter(filter interface{}) (models.Coach, error)
	GetCoachPage(paginated *utils.PaginationRequest) (interface{}, error)
	GetCoachCollection() *mongo.Collection

	AddPlayer(player *models.Player) (string, error)
	DeletePlayer(id primitive.ObjectID) error
	UpdatePlayer(coach *models.Player) error
	GetPlayerById(id primitive.ObjectID) (models.Player, error)
	GetPlayerByFilter(filter interface{}) (models.Player, error)
	GetPlayerPage(paginated *utils.PaginationRequest) (interface{}, error)
	GetPlayerCollection() *mongo.Collection

	AddTeam(team *models.Team) (string, error)
	DeleteTeam(id primitive.ObjectID) error
	UpdateTeam(team *models.Team) error
	GetTeamById(id primitive.ObjectID) (models.Team, error)
	GetTeamPage(paginated *utils.PaginationRequest) (interface{}, error)

	GetUserById(id string, collection *mongo.Collection) (interface{}, error)
}

type Client struct {
	DB *mongo.Database

	coachCollection  *mongo.Collection
	playerCollection *mongo.Collection
	teamCollection   *mongo.Collection
}

func (c *Client) Connect(ctx *context.Context) (interface{}, error) {
	var dbName string

	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("URI")))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = client.Connect(*ctx)
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
		coachCollection:  client.Database(dbName).Collection("Coaches"),
		playerCollection: client.Database(dbName).Collection("Players"),
		teamCollection:   client.Database(dbName).Collection("Teams"),
	}

	return c.DB.Client(), nil
}

func (c *Client) Disconnect(ctx *context.Context) error {
	return c.DB.Client().Disconnect(*ctx)
}

func (c *Client) GetUserById(id string, collection *mongo.Collection) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	newId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user interface{}
	filter := bson.D{{Key: "_id", Value: newId}}
	err = collection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		return nil, err
	}

	if collection.Name() == "Players" {
		return user.(models.Player), nil
	} else {
		return user.(models.Coach), nil
	}
}
