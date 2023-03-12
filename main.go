package main

import (
	"context"
	"football-manager-go/middleware"
	"football-manager-go/routes"
	"football-manager-go/utils/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func main() {
	err := godotenv.Load("example.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	client := &db.Client{}
	if _, err := client.Connect(&ctx); err != nil {
		log.Fatal("Cannot connect to database: ", err.Error())
	}

	defer func() {
		_ = client.Disconnect(&ctx)
	}()

	router := gin.Default()
	router.Use(middleware.InstallDatabase(client))
	routes.DeclareRoutes(router)

	log.Fatal(router.Run(":" + os.Getenv("PORT")))
}
