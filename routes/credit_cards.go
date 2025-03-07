package routes

import (
	"go_income_outflow/middleware"
	"go_income_outflow/pkg/controller"
	"go_income_outflow/pkg/repository"
	"go_income_outflow/pkg/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func creditCardRoutes(r *gin.Engine, db *gorm.DB) {
	creditCardRepo := repository.NewCreditCardRepository(db)
	creditCardService := service.NewCreditCardService(creditCardRepo)
	creditCardController := controller.NewCreditCardController(creditCardService)

	creditCardGroup := r.Group("/credit_cards")
	creditCardGroup.Use(middleware.Auth)
	creditCardGroup.GET("", creditCardController.Index)
	creditCardGroup.POST("", creditCardController.Store)
	creditCardGroup.GET("/:id", creditCardController.Show)

	creditCardGroup.PUT("/:id", creditCardController.Update)
	creditCardGroup.PATCH("/:id", creditCardController.Update)

	creditCardGroup.DELETE("/:id", creditCardController.Destroy)
}
