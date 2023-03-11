package main

import (
	"football-manager-go/controllers"
	"football-manager-go/middleware"
	"football-manager-go/routes"
	"football-manager-go/utils/db"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	client := &db.Client{}
	if _, err := client.Connect(); err != nil {
		log.Fatal("Cannot connect to database: ", err.Error())
	}

	defer func() {
		_ = client.Disconnect()
	}()

	router := gin.Default()
	v1 := router.Group(routes.Base)

	router.Use(middleware.InstallDatabase(client))

	coachRoutes := v1.Group("/coaches")
	{
		coachRoutes.GET("/:id", controllers.GetCoach())
	}

	log.Fatal(router.Run(":" + os.Getenv("PORT")))
}
