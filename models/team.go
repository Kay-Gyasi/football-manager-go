package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Team struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name            string             `json:"name" bson:"name" validate:"required"`
	DateEstablished primitive.DateTime `json:"dateEstablished" bson:"date_established"`
	StadiumName     string             `json:"stadiumName" bson:"stadium_name"`
}
