package middleware

import (
	"fmt"
	"football-manager-go/utils"
	"football-manager-go/utils/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
)

func RequiresAuth(c *gin.Context) {
	database := c.MustGet("database").(db.IDatabase)
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		utils.FailureResponse(c, http.StatusUnauthorized, "Authorization Header missing")
		c.Abort()
		return
	}

	results := strings.Split(authHeader, " ")
	if len(results) != 2 {
		utils.FailureResponse(c, http.StatusUnauthorized, "Authorization Header must start with `Bearer`")
		c.Abort()
		return
	}

	// TODO: Validate the token

	claims, err := utils.GetClaimsFromToken(results[1])
	if err != nil {
		fmt.Println(err.Error())
		utils.FailureResponse(c, http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}

	if len(claims.Roles) < 1 {
		utils.FailureResponse(c, 400, "Invalid user")
		c.Abort()
		return
	}

	var dbCollection *mongo.Collection
	for _, role := range claims.Roles {
		if role != "player" || role != "coach" || role != "admin" {
			utils.FailureResponse(c, 400, "Invalid role")
			c.Abort()
			return
		}

		if role == "player" {
			dbCollection = database.GetPlayerCollection()
		}
		if role == "coach" {
			dbCollection = database.GetCoachCollection()
		}
	}

	currentUser, err := database.GetUserById(claims.ID, dbCollection)
	if err != nil {
		utils.FailureResponse(c, 500, "Something went wrong")
		c.Abort()
		return
	}

	c.Set("currentUser", currentUser)
}
