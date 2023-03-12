package controllers

import (
	"errors"
	"football-manager-go/models"
	"football-manager-go/models/embedded"
	"football-manager-go/utils"
	"football-manager-go/utils/db"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
)

var Validate = validator.New()

func GetCoach(c *gin.Context) {
	id := c.Param("id")
	database, ok := c.MustGet("database").(db.IDatabase)
	if !ok {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		c.Abort()
		return
	}

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id not provided"})
		c.Abort()
		return
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		c.Abort()
		return
	}

	coach, err := database.GetCoachById(objectId)
	if err != nil {
		if errors.Is(err, utils.InvalidObjectID) || errors.Is(err, utils.CantFindCoach) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			c.Abort()
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, coach)
	c.Done()
}

func InsertCoach(c *gin.Context) {
	var coach models.Coach
	dbClient := c.MustGet("database").(db.IDatabase)

	if err := c.BindJSON(&coach); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	if len(coach.Password) < 1 {
		utils.FailureResponse(c, 400, "Invalid password")
		c.Abort()
		return
	}

	validationErr := Validate.Struct(coach)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		c.Abort()
		return
	}

	if coach.Type != embedded.IsCoach {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user type"})
		c.Abort()
		return
	}

	var usernameBuilder strings.Builder
	usernameBuilder.WriteString(coach.Firstname)
	usernameBuilder.WriteString(" ")
	usernameBuilder.WriteString(coach.Lastname)
	coach.Username = usernameBuilder.String()

	coach.PasswordHash = HashPassword(coach.Password)

	// checking if coach already exists in database
	filter := bson.M{"user.firstname": coach.Firstname, "user.lastname": coach.Lastname,
		"user.email": coach.Email, "user.phone": coach.Phone}

	filtered, err := dbClient.GetCoachByFilter(filter)
	if filtered.ID != primitive.NilObjectID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Coach already exists"})
		c.Abort()
		return
	}

	var id string
	id, err = dbClient.AddCoach(&coach)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, id)
	c.Done()
}

func UpdateCoach(c *gin.Context) {
	var coach models.Coach
	dbClient := c.MustGet("database").(db.IDatabase)

	if err := c.BindJSON(&coach); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	validationErr := Validate.Struct(coach)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	if coach.ID == primitive.NilObjectID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		c.Abort()
		return
	}

	if coach.Type != embedded.IsCoach {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user type"})
		c.Abort()
		return
	}

	err := dbClient.UpdateCoach(&coach)
	if err != nil {
		if errors.Is(err, utils.InvalidObjectID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			c.Abort()
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusNoContent, "")
	c.Done()
}

func DeleteCoach(c *gin.Context) {
	id := c.Param("id")
	database, ok := c.MustGet("database").(db.IDatabase)
	if !ok {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		c.Abort()
		return
	}

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id not provided"})
		c.Abort()
		return
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		c.Abort()
		return
	}

	err = database.DeleteCoach(objectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusNoContent, "")
	c.Done()
}

func GetCoachPage(c *gin.Context) {
	var paginated utils.PaginationRequest
	dbClient := c.MustGet("database").(db.IDatabase)

	if err := c.BindJSON(&paginated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	if err := Validate.Struct(paginated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	if paginated.PageNumber <= 0 {
		paginated.PageNumber = 1
	}
	if paginated.PageSize <= 0 {
		paginated.PageSize = 10
	}

	page, err := dbClient.GetCoachPage(&paginated)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, page.(utils.PaginationResponse))
	c.Done()
}
