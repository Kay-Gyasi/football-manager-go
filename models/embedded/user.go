package embedded

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserType int
type UserRole int

const (
	IsPlayer UserType = iota
	IsCoach
	IsDeveloper
)

const (
	Player UserRole = iota
	Coach
	Admin
)

type User struct {
	Username     string             `json:"username" bson:"username"`
	Firstname    string             `json:"firstname" bson:"firstname" validate:"required,min=2, max=30"`
	Lastname     string             `json:"lastname" bson:"lastname" validate:"required,min=2, max=30"`
	Type         UserType           `json:"type" bson:"type" validate:"required"`
	DateOfBirth  primitive.DateTime `json:"dateOfBirth" bson:"date_of_birth"`
	Email        string             `json:"email" bson:"email"`
	PasswordHash string             `json:"-" bson:"password_hash"`
	Phone        string             `json:"phone" bson:"phone"`
	Roles        []Role             `json:"-" bson:"roles"`
}

type Role struct {
	Name UserRole `json:"-" bson:"name"`
}
