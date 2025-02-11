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
	transactionCategoryGroup.GET("", middleware.Auth, transactionCategoryController.Index)
	transactionCategoryGroup.POST("", middleware.Auth, transactionCategoryController.Store)
	transactionCategoryGroup.GET("/:id", middleware.Auth, transactionCategoryController.Show)

	transactionCategoryGroup.PUT("/:id", middleware.Auth, transactionCategoryController.Update)
	transactionCategoryGroup.PATCH("/:id", middleware.Auth, transactionCategoryController.Update)

	transactionCategoryGroup.DELETE("/:id", middleware.Auth, transactionCategoryController.Destroy)
}
