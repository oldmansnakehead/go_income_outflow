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
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)

	userRepo := repository.NewUserRepository(db, refreshTokenRepo)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	userGroup := r.Group("/users")
	userGroup.GET("/refresh_token", userController.RefreshToken)

	userGroup.POST("", userController.Store)
	userGroup.POST("/login", userController.Login)
	userGroup.POST("/logout", middleware.Auth, userController.Logout)
	userGroup.POST("/test-auth", middleware.Auth, userController.TestAuth)
}
