package routes

import (
	"go_income_outflow/middleware"
	"go_income_outflow/pkg/controller"
	"go_income_outflow/pkg/repository"
	"go_income_outflow/pkg/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func transactionRoutes(r *gin.Engine, db *gorm.DB) {
	creditCardRepo := repository.NewCreditCardRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	creditCardDebtRepo := repository.NewCreditCardDebtRepository(db, creditCardRepo)
	transactionCategoryRepo := repository.NewTransactionCategoryRepository(db)

	transactionRepo := repository.NewTransactionRepository(db, accountRepo, creditCardRepo, creditCardDebtRepo, transactionCategoryRepo)
	transactionService := service.NewTransactionService(transactionRepo)
	transactionController := controller.NewTransactionController(transactionService)

	transactionGroup := r.Group("/transactions")
	transactionGroup.Use(middleware.Auth)
	transactionGroup.GET("", transactionController.Index)
	transactionGroup.POST("", transactionController.Store)
	transactionGroup.GET("/:id", transactionController.Show)

	transactionGroup.DELETE("/:id", transactionController.Destroy)
}
