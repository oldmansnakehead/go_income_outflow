package routes

import (
	"go_income_outflow/pkg/controller"

	"github.com/gin-gonic/gin"
)

func userRoutes(r *gin.Engine) {
	userController := controller.User{}
	userGroup := r.Group("/users")
	userGroup.POST("", userController.Create)
}
