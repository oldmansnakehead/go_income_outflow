package controller

import (
	"go_income_outflow/pkg/custom/request"
	"go_income_outflow/pkg/model"
	"go_income_outflow/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	CreditCardDebtControllerUseCase interface {
		Store(ctx *gin.Context)
	}

	creditCardDebtController struct {
		service service.CreditCardDebtServiceUseCase
	}
)

func NewCreditCardDebtController(service service.CreditCardDebtServiceUseCase) CreditCardDebtControllerUseCase {
	return &creditCardDebtController{service: service}
}

func (c *creditCardDebtController) Store(ctx *gin.Context) {
	var form model.CreditCardDebtRequest
	validateCtx := request.NewCustomRequest(ctx)
	if err := validateCtx.BindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	debts, err := c.service.CreateDebt(&form)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "debt created successfully",
		"data":    debts,
	})
}
