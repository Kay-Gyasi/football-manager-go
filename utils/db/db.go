package db

import (
	"context"
	"football-manager-go/models"
	"football-manager-go/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

type IDatabase interface {
	Connect() (interface{}, error)
	Disconnect() error

	AddCoach(coach models.Coach) (string, error)
	DeleteCoach(id string) error
	UpdateCoach(coach models.Coach) error
	GetCoachById(id string) (models.Coach, error)
	GetCoachPage(paginated utils.PaginationRequest) (interface{}, error)

	AddPlayer(player models.Player) (string, error)
	DeletePlayer(id string) error
	UpdatePlayer(coach models.Player) error
	GetPlayerById(id string) (models.Player, error)
	GetPlayerPage(paginated utils.PaginationRequest) (interface{}, error)

	AddTeam(team models.Team) (string, error)
	DeleteTeam(id string) error
	UpdateTeam(team models.Team) error
	GetTeamById(id string) (models.Team, error)
	GetTeamPage(paginated utils.PaginationRequest) (interface{}, error)
}

type Client struct {
	DB *mongo.Database

	coachCollection  *mongo.Collection
	playerCollection *mongo.Collection
	teamCollection   *mongo.Collection
}

func (c *Client) Connect() (interface{}, error) {
	var dbName string

	client, err := mongo.Connect(context.Background())
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
	return c.DB.Client().Disconnect(context.Background())
}
