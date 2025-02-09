package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitialRoute(r *gin.Engine, db *gorm.DB) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "welcome",
		})
	})

	userRoutes(r, db)
	accountRoutes(r, db)
	creditCardRoutes(r, db)
}
