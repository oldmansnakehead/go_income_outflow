package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitialRoute(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "welcome",
		})
	})

	userRoutes(r)
}
