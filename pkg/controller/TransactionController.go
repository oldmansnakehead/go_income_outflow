package controller

import (
	"fmt"
	"go_income_outflow/entities"
	"go_income_outflow/helpers"
	"go_income_outflow/pkg/custom/request"
	"go_income_outflow/pkg/model"
	"go_income_outflow/pkg/service"
	"net/http"
	"strconv"

	"go_income_outflow/pkg/controller/common"
	modelCommon "go_income_outflow/pkg/model/common"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type (
	TransactionControllerUseCase interface {
		common.ControllerUseCase
	}
	transactionController struct {
		service service.TransactionServiceUseCase
	}
)

func NewTransactionController(service service.TransactionServiceUseCase) TransactionControllerUseCase {
	return &transactionController{service: service}
}

func (c *transactionController) Index(ctx *gin.Context) {
	filters := helpers.ParseQueryString(ctx)

	transactions, err := c.service.GetWithFilters(filters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}

func (c *transactionController) Store(ctx *gin.Context) {
	var form model.TransactionRequest
	validateCtx := request.NewCustomRequest(ctx)
	if err := validateCtx.BindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction := entities.Transaction{}
	copier.Copy(&transaction, form)

	if err := c.service.CreateTransaction(&transaction, form); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "transaction created successfully",
		"data":    transaction,
	})
}

func (c *transactionController) Show(ctx *gin.Context) {
	id := ctx.Param("id")

	form := modelCommon.CommonRequest{}

	helpers.ParseQueryString(ctx)

	withArray := ctx.QueryArray("with[]")
	if len(withArray) > 0 {
		form.With = withArray
	}

	var transaction entities.Transaction

	err := c.service.FirstWithRelations(&transaction, form.With)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("transaction with ID %s not found", id),
		})
		return
	}

	ctx.JSON(http.StatusOK, (&model.Transaction{}).Response(&transaction))
}

func (c *transactionController) Update(ctx *gin.Context) {

}

func (c *transactionController) Destroy(ctx *gin.Context) {
	id := ctx.Param("id")

	uintID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}
	transaction := entities.Transaction{
		Model: gorm.Model{ID: uint(uintID)},
	}

	if err := c.service.DeleteTransaction(&transaction); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "transaction deleted successfully",
	})
}
