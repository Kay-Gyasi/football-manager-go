package models

import (
	"football-manager-go/models/embedded"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Position int
type Nationality int

const (
	GoalKeeper Position = iota
	Defender
	Midfielder
	Forward
)

const (
	Ghanaian Nationality = iota
	Swiss
	Nigerian
	Argentinian
	English
)

type Player struct {
	embedded.User
	ID                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	TeamID            primitive.ObjectID `json:"teamId" bson:"team_id"`
	JerseyName        string             `json:"jerseyName" bson:"jersey_name"`
	JerseyNumber      uint8              `json:"jerseyNumber" bson:"jersey_number"`
	PrimaryPosition   Position           `json:"primaryPosition" bson:"primary_position"`
	SecondaryPosition Position           `json:"secondaryPosition" bson:"secondary_position"`
	Nationality       Nationality        `json:"nationality" bson:"nationality"`
}
