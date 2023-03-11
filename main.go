package main

import (
	"context"
	"football-manager-go/controllers"
	"football-manager-go/middleware"
	"football-manager-go/routes"
	"football-manager-go/utils/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	client := &db.Client{}
	if _, err := client.Connect(&ctx); err != nil {
		log.Fatal("Cannot connect to database: ", err.Error())
	}

	defer func() {
		_ = client.Disconnect(&ctx)
	}()

	err := godotenv.Load("example.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	router := gin.Default()
	router.Use(middleware.InstallDatabase(client))
	v1 := router.Group(routes.Base)

	coachRoutes := v1.Group("/coaches")
	{
		coachRoutes.GET("/:id", controllers.GetCoach)
	}

	log.Fatal(router.Run(":" + os.Getenv("PORT")))
}
