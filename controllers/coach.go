package controllers

import (
	"errors"
	"football-manager-go/utils"
	"football-manager-go/utils/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

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

	coach, err := database.GetCoachById(id)
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
}
