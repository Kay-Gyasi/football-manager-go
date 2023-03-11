package db

import (
	"context"
	"football-manager-go/models"
	"football-manager-go/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *Client) AddPlayer(player models.Player) (string, error) {
	response, err := c.playerCollection.InsertOne(context.Background(), player)
	if err != nil {
		return "", err
	}

	return response.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (c *Client) DeletePlayer(id string) error {
	newId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: newId}}
	_, err = c.playerCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdatePlayer(player models.Player) error {
	if player.ID == primitive.NilObjectID {
		return utils.CantUpdateNewDocument
	}

	filter := bson.D{{Key: "_id", Value: player.ID}}
	update := bson.M{"$set": bson.M{"username": player.Username, "firstname": player.Firstname,
		"lastname": player.Lastname, "type": player.Type, "date_of_birth": player.DateOfBirth,
		"email": player.Email, "phone": player.Phone, "jersey_name": player.JerseyName,
		"jersey_number": player.JerseyNumber, "primary_position": player.PrimaryPosition,
		"secondary_position": player.SecondaryPosition, "nationality": player.Nationality}}
	_, err := c.playerCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetPlayerById(id string) (models.Player, error) {
	var player models.Player
	newId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return player, utils.InvalidObjectID
	}

	err = c.playerCollection.FindOne(context.Background(),
		bson.D{{"_id", newId}}).Decode(&player)
	if err != nil {
		return player, err
	}

	return player, nil
}

func (c *Client) GetPlayerPage(paginated utils.PaginationRequest) (interface{}, error) {
	return nil, nil
}
