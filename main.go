package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	router := gin.Default()

	log.Fatal(router.Run())
}
