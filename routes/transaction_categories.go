package routes

import (
	"go_income_outflow/middleware"
	"go_income_outflow/pkg/controller"
	"go_income_outflow/pkg/repository"
	"go_income_outflow/pkg/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func transactionCategoryRoutes(r *gin.Engine, db *gorm.DB) {
	transactionCategoryRepo := repository.NewTransactionCategoryRepository(db)
	transactionCategoryService := service.NewTransactionCategoryService(transactionCategoryRepo)
	transactionCategoryController := controller.NewTransactionCategoryController(transactionCategoryService)

	transactionCategoryGroup := r.Group("/transaction_categories")
	transactionCategoryGroup.Use(middleware.Auth)
	transactionCategoryGroup.GET("", transactionCategoryController.Index)
	transactionCategoryGroup.POST("", transactionCategoryController.Store)
	transactionCategoryGroup.GET("/:id", transactionCategoryController.Show)

	transactionCategoryGroup.PUT("/:id", transactionCategoryController.Update)
	transactionCategoryGroup.PATCH("/:id", transactionCategoryController.Update)

	transactionCategoryGroup.DELETE("/:id", transactionCategoryController.Destroy)
}
