package db

import (
	"context"
	"football-manager-go/models"
	"football-manager-go/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *Client) AddTeam(team models.Team) (string, error) {
	response, err := c.teamCollection.InsertOne(context.Background(), team)
	if err != nil {
		return "", err
	}

	return response.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (c *Client) DeleteTeam(id string) error {
	newId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: newId}}
	_, err = c.teamCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateTeam(team models.Team) error {
	if team.ID == primitive.NilObjectID {
		return utils.CantUpdateNewDocument
	}

	filter := bson.D{{Key: "_id", Value: team.ID}}
	update := bson.M{"$set": bson.M{"name": team.Name, "date_established": team.DateEstablished,
		"stadium_name": team.StadiumName}}
	_, err := c.teamCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetTeamById(id string) (models.Team, error) {
	var team models.Team
	err := c.teamCollection.FindOne(context.Background(),
		bson.D{{"_id", primitive.ObjectIDFromHex(id)}}).Decode(&team)
	if err != nil {
		return team, err
	}

	return team, nil
}

func (c *Client) GetTeamPage(paginated utils.PaginationRequest) (interface{}, error) {
	return nil, nil
}
