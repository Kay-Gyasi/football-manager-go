package middleware

import (
	"football-manager-go/utils/db"
	"github.com/gin-gonic/gin"
)

func InstallDatabase(database db.IDatabase) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("database", database)
		c.Next()
	}
}
