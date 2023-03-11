package controllers

import (
	"errors"
	"football-manager-go/utils"
	"football-manager-go/utils/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetCoach() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		database, ok := c.MustGet("database").(db.IDatabase)
		if !ok {
			c.JSON(500, "Internal Server Error")
			c.Abort()
			return
		}

		if id == "" {
			c.JSON(http.StatusBadRequest, "Id not provided")
			c.Abort()
			return
		}

		coach, err := database.GetCoachById(id)
		if err != nil {
			if errors.Is(err, utils.InvalidObjectID) || errors.Is(err, utils.CantFindCoach) {
				c.JSON(http.StatusBadRequest, "Invalid id")
				c.Abort()
				return
			}

			c.JSON(http.StatusInternalServerError, "Something went wrong")
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, coach)
	}
}
