package routes

import (
	"football-manager-go/controllers"
	"github.com/gin-gonic/gin"
)

const Base string = "/api/v1"

func DeclareRoutes(router *gin.Engine) {
	v1 := router.Group(Base)

	coachRoutes := v1.Group("/coaches")
	{
		coachRoutes.GET("/:id", controllers.GetCoach)
		coachRoutes.POST("/page", controllers.GetCoachPage)
		coachRoutes.POST("", controllers.InsertCoach)
		coachRoutes.PUT("", controllers.UpdateCoach)
		coachRoutes.DELETE("/:id", controllers.DeleteCoach)
	}

	playerRoutes := v1.Group("/players")
	{
		playerRoutes.GET("/:id", controllers.GetPlayer)
		playerRoutes.POST("/page", controllers.GetPlayerPage)
		playerRoutes.POST("", controllers.InsertPlayer)
		playerRoutes.PUT("", controllers.UpdatePlayer)
		playerRoutes.DELETE("/:id", controllers.DeletePlayer)
	}
}
