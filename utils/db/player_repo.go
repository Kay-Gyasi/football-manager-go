package db

import (
	"context"
	"football-manager-go/models"
	"football-manager-go/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	options2 "go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"time"
)

func (c *Client) AddPlayer(player *models.Player) (string, error) {
	response, err := c.playerCollection.InsertOne(context.Background(), player)
	if err != nil {
		return "", err
	}

	return response.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (c *Client) DeletePlayer(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: id}}
	_, err := c.playerCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdatePlayer(player *models.Player) error {
	if player.ID == primitive.NilObjectID {
		return utils.CantUpdateNewDocument
	}

	filter := bson.D{{Key: "_id", Value: player.ID}}
	update := bson.M{"$set": bson.M{"user.username": player.Username, "user.firstname": player.Firstname,
		"user.lastname": player.Lastname, "user.type": player.Type, "user.date_of_birth": player.DateOfBirth,
		"user.email": player.Email, "user.phone": player.Phone, "jersey_name": player.JerseyName,
		"jersey_number": player.JerseyNumber, "primary_position": player.PrimaryPosition,
		"secondary_position": player.SecondaryPosition, "nationality": player.Nationality}}
	_, err := c.playerCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetPlayerById(id primitive.ObjectID) (models.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var player models.Player

	err := c.playerCollection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&player)
	if err != nil {
		return player, utils.CantFindPlayer
	}

	return player, nil
}

func (c *Client) GetPlayerByFilter(filter interface{}) (models.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var player models.Player

	err := c.playerCollection.FindOne(ctx, filter).Decode(&player)
	if err != nil {
		return player, utils.CantFindPlayer
	}

	return player, nil
}

func (c *Client) GetPlayerPage(paginated *utils.PaginationRequest) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var page []models.Player

	options := options2.Find()
	options.SetSkip(int64((paginated.PageNumber - 1) * paginated.PageSize))
	options.SetLimit(int64(paginated.PageSize))

	cursor, err := c.playerCollection.Find(ctx, bson.D{}, options)
	if err != nil {
		return nil, err
	}

	// Iterate over the results and add them to the pages array
	for cursor.Next(ctx) {
		var player models.Player
		err := cursor.Decode(&player)
		if err != nil {
			return nil, err
		}

		page = append(page, player)
	}

	var totalItems int64
	// Calculate total items and total pages
	totalItems, err = c.playerCollection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	totalPages := int(math.Ceil(float64(totalItems) / float64(paginated.PageSize)))

	// Create the pagination response
	paginationRes := utils.PaginationResponse{
		Data:       page,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
		Page:       paginated.PageNumber,
		PageSize:   paginated.PageSize,
	}

	defer cursor.Close(ctx)
	return paginationRes, nil
}

func (c *Client) GetPlayerCollection() *mongo.Collection {
	return c.playerCollection
}
