package db

import (
	"context"
	"football-manager-go/models"
	"football-manager-go/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *Client) AddCoach(coach models.Coach) (string, error) {
	response, err := c.coachCollection.InsertOne(context.Background(), coach)
	if err != nil {
		return "", err
	}

	return response.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (c *Client) DeleteCoach(id string) error {
	newId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: newId}}
	_, err = c.coachCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateCoach(coach models.Coach) error {
	if coach.ID == primitive.NilObjectID {
		return utils.CantUpdateNewDocument
	}

	filter := bson.D{{Key: "_id", Value: coach.ID}}
	update := bson.M{"$set": bson.M{"username": coach.Username, "firstname": coach.Firstname,
		"lastname": coach.Lastname, "type": coach.Type, "date_of_birth": coach.DateOfBirth,
		"email": coach.Email, "phone": coach.Phone, "years_of_experience": coach.YearsOfExperience}}
	_, err := c.coachCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetCoachById(id string) (models.Coach, error) {
	var coach models.Coach
	newId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return coach, utils.InvalidObjectID
	}

	err = c.coachCollection.FindOne(context.Background(),
		bson.D{{"_id", newId}}).Decode(&coach)
	if err != nil {
		return coach, utils.CantFindCoach
	}

	return coach, nil
}

func (c *Client) GetCoachPage(paginated utils.PaginationRequest) (interface{}, error) {
	return nil, nil
}
