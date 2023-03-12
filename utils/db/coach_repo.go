package db

import (
	"context"
	"football-manager-go/models"
	"football-manager-go/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	options2 "go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"time"
)

func (c *Client) AddCoach(coach *models.Coach) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	response, err := c.coachCollection.InsertOne(ctx, coach)
	if err != nil {
		return "", err
	}

	return response.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (c *Client) DeleteCoach(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: id}}
	_, err := c.coachCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateCoach(coach *models.Coach) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if coach.ID == primitive.NilObjectID {
		return utils.InvalidObjectID
	}

	filter := bson.D{{Key: "_id", Value: coach.ID}}
	update := bson.M{"$set": bson.M{"user.username": coach.Username, "user.firstname": coach.Firstname,
		"user.lastname": coach.Lastname, "user.type": coach.Type, "user.date_of_birth": coach.DateOfBirth,
		"user.email": coach.Email, "user.phone": coach.Phone, "years_of_experience": coach.YearsOfExperience}}
	_, err := c.coachCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetCoachById(id primitive.ObjectID) (models.Coach, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var coach models.Coach

	err := c.coachCollection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&coach)
	if err != nil {
		return coach, utils.CantFindCoach
	}

	return coach, nil
}

func (c *Client) GetCoachByFilter(filter interface{}) (models.Coach, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var coach models.Coach

	err := c.coachCollection.FindOne(ctx, filter).Decode(&coach)
	if err != nil {
		return coach, utils.CantFindCoach
	}

	return coach, nil
}

func (c *Client) GetCoachPage(paginated *utils.PaginationRequest) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var page []models.Coach

	options := options2.Find()
	options.SetSkip(int64((paginated.PageNumber - 1) * paginated.PageSize))
	options.SetLimit(int64(paginated.PageSize))

	cursor, err := c.coachCollection.Find(ctx, bson.D{}, options)
	if err != nil {
		return nil, err
	}

	// Iterate over the results and add them to the pages array
	for cursor.Next(ctx) {
		var coach models.Coach
		err := cursor.Decode(&coach)
		if err != nil {
			return nil, err
		}

		page = append(page, coach)
	}

	var totalItems int64
	// Calculate total items and total pages
	totalItems, err = c.coachCollection.CountDocuments(ctx, bson.D{})
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
