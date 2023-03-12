package controllers

import (
	"football-manager-go/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// TODO: Work on this, modify jwt.go, authorize routes, CORS, health checks, rate limiting, swagger integration

type LoginCommand struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func Login(c *gin.Context) {
	var loginCommand LoginCommand
	//dbClient := c.MustGet("database").(db.IDatabase)

	if err := c.BindJSON(&loginCommand); err != nil {
		utils.FailureResponse(c, 400, "Invalid input")
		c.Abort()
		return
	}

	// write method to get player/coach by email

	//passwordHash := VerifyPassword(loginCommand.Password)
	c.Done()
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := ""

	if err != nil {
		msg = "Login or Password is incorrect"
		valid = false
	}
	return valid, msg
}