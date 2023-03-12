package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Payload map[string]interface{}

func SuccessResponse(c *gin.Context, statusCode int, message string, payload Payload) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"message": message,
		"status":  http.StatusText(statusCode),
		"payload": payload,
	})
}

func FailureResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"error":   message,
		"status":  http.StatusText(statusCode),
	})
}
