package models

import (
	"football-manager-go/models/embedded"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Coach struct {
	embedded.User
	ID                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	TeamID            primitive.ObjectID `json:"teamId" bson:"team_id"`
	YearsOfExperience uint8              `json:"yearsOfExperience" bson:"years_of_experience"`
	IsMain            bool               `bson:"is_main" json:"isMain"`
}
