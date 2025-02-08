package controller

// controller = handler

import (
	"fmt"
	"go_income_outflow/entities"
	"go_income_outflow/pkg/model"
	"go_income_outflow/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	service *service.AccountService
}

func NewAccountController(service *service.AccountService) *AccountController {
	return &AccountController{service: service}
}

func (ac *AccountController) Index(ctx *gin.Context) {
}

func (ac *AccountController) Store(ctx *gin.Context) {
	var form model.AccountRequest
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account := entities.Account{
		Name:   form.Name,
		UserID: form.UserID,
	}

	if err := ac.service.CreateAccountWithRelations(&account, []string{"User"}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create account",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Account created successfully",
		"data":    (&model.Account{}).EntitiesToModel(&account).ToResponse(), // keyword = struct literal
	})
}

func (ac *AccountController) Show(ctx *gin.Context) {
}

func (ac *AccountController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	fmt.Println(id)
}

func (ac *AccountController) Destroy(ctx *gin.Context) {
}
