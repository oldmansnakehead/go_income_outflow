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
	accountGroup.Use(middleware.Auth)
	accountGroup.GET("", accountController.Index)
	accountGroup.POST("", accountController.Store)
	accountGroup.GET("/:id", accountController.Show)

	accountGroup.PUT("/:id", accountController.Update)
	accountGroup.PATCH("/:id", accountController.Update)

	accountGroup.DELETE("/:id", accountController.Destroy)

	accountGroup.GET("/currencies", accountController.GetCurrencies)
	accountGroup.GET("/total_amount", accountController.GetTotalAmount)
}
