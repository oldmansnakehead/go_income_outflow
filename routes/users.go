package routes

import (
	"go_income_outflow/middleware"
	"go_income_outflow/pkg/controller"
	"go_income_outflow/pkg/repository"
	"go_income_outflow/pkg/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func userRoutes(r *gin.Engine, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	userGroup := r.Group("/users")
	userGroup.POST("", userController.Store)
	userGroup.POST("/login", middleware.Auth, userController.Login)
	userGroup.POST("/test-auth", middleware.Auth, userController.TestAuth)
}
