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
	creditCardGroup.GET("", middleware.Auth, creditCardController.Index)
	creditCardGroup.POST("", middleware.Auth, creditCardController.Store)
	creditCardGroup.GET("/:id", middleware.Auth, creditCardController.Show)

	creditCardGroup.PUT("/:id", middleware.Auth, creditCardController.Update)
	creditCardGroup.PATCH("/:id", middleware.Auth, creditCardController.Update)

	creditCardGroup.DELETE("/:id", middleware.Auth, creditCardController.Destroy)
}
