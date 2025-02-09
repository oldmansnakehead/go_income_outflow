package routes

import (
	"go_income_outflow/middleware"
	"go_income_outflow/pkg/controller"
	"go_income_outflow/pkg/repository"
	"go_income_outflow/pkg/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func accountRoutes(r *gin.Engine, db *gorm.DB) {
	accountRepo := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepo)
	accountController := controller.NewAccountController(accountService)

	accountGroup := r.Group("/accounts")
	accountGroup.GET("", middleware.Auth, accountController.Index)
	accountGroup.POST("", middleware.Auth, accountController.Store)
	accountGroup.GET("/:id", middleware.Auth, accountController.Show)

	accountGroup.PUT("/:id", middleware.Auth, accountController.Update)
	accountGroup.PATCH("/:id", middleware.Auth, accountController.Update)

	accountGroup.DELETE("/:id", middleware.Auth, accountController.Destroy)
}
