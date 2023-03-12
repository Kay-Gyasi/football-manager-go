package controllers

import (
	"errors"
	"football-manager-go/models"
	"football-manager-go/models/embedded"
	"football-manager-go/utils"
	"football-manager-go/utils/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
)

func GetPlayer(c *gin.Context) {
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

	coach, err := database.GetPlayerById(objectId)
	if err != nil {
		if errors.Is(err, utils.InvalidObjectID) || errors.Is(err, utils.CantFindPlayer) {
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

func InsertPlayer(c *gin.Context) {
	var player models.Player
	dbClient := c.MustGet("database").(db.IDatabase)

	if err := c.BindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	validationErr := Validate.Struct(player)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	if player.Type != embedded.IsPlayer {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user type"})
		c.Abort()
		return
	}

	var usernameBuilder strings.Builder
	usernameBuilder.WriteString(player.Firstname)
	usernameBuilder.WriteString(" ")
	usernameBuilder.WriteString(player.Lastname)
	player.Username = usernameBuilder.String()

	// checking if player already exists in database
	filter := bson.M{"user.firstname": player.Firstname, "user.lastname": player.Lastname,
		"user.email": player.Email, "user.phone": player.Phone}

	filtered, err := dbClient.GetPlayerByFilter(filter)
	if filtered.ID != primitive.NilObjectID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Coach already exists"})
		c.Abort()
		return
	}

	var id string
	id, err = dbClient.AddPlayer(&player)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, id)
	c.Done()
}

func UpdatePlayer(c *gin.Context) {
	var player models.Player
	dbClient := c.MustGet("database").(db.IDatabase)

	if err := c.BindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	validationErr := Validate.Struct(player)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	if player.ID == primitive.NilObjectID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		c.Abort()
		return
	}

	if player.Type != embedded.IsPlayer {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user type"})
		c.Abort()
		return
	}

	err := dbClient.UpdatePlayer(&player)
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

func DeletePlayer(c *gin.Context) {
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

	err = database.DeletePlayer(objectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusNoContent, "")
	c.Done()
}

func GetPlayerPage(c *gin.Context) {
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

	page, err := dbClient.GetPlayerPage(&paginated)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, page.(utils.PaginationResponse))
	c.Done()
}
