package routes

import (
	"go_income_outflow/middleware"
	"go_income_outflow/pkg/controller"
	"go_income_outflow/pkg/repository"
	"go_income_outflow/pkg/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func creditCardDebtRoutes(r *gin.Engine, db *gorm.DB) {
	creditCardRepo := repository.NewCreditCardRepository(db)

	creditCardDebtRepo := repository.NewCreditCardDebtRepository(db, creditCardRepo)
	creditCardDebtService := service.NewCreditCardDebtService(creditCardDebtRepo)
	creditCardDebtController := controller.NewCreditCardDebtController(creditCardDebtService)

	creditCardDebtGroup := r.Group("/credit_card_debts")
	// creditCardDebtGroup.GET("", middleware.Auth, creditCardDebtController.Index)
	creditCardDebtGroup.POST("", middleware.Auth, creditCardDebtController.Store)
	// creditCardDebtGroup.GET("/:id", middleware.Auth, creditCardDebtController.Show)

	// creditCardDebtGroup.DELETE("/:id", middleware.Auth, creditCardDebtController.Destroy)
}
