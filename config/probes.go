package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "i'm healthy",
	})
}

func Ready(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "i'm ready",
	})
}
