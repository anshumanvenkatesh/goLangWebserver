package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// First experimentation with routes in a separate package. 404 having its own file was definitely not intended!

func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "Umm... Maybe you were aiming for /aggregate or /buildings? That's all for the show now.",
		})
	}
}
