package controller

// controller = handler

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

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	TransactionCategoryControllerUseCase interface {
		common.ControllerUseCase
	}

	transactionCategoryController struct {
		service service.TransactionCategoryServiceUseCase
	}
)

func NewTransactionCategoryController(service service.TransactionCategoryServiceUseCase) TransactionCategoryControllerUseCase {
	return &transactionCategoryController{service: service}
}

func (c *transactionCategoryController) Index(ctx *gin.Context) {
	filters := helpers.ParseQueryString(ctx)

	items, err := c.service.GetWithFilters(filters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, items)
}

func (c *transactionCategoryController) Store(ctx *gin.Context) {
	var form model.TransactionCategoryRequest
	validateCtx := request.NewCustomRequest(ctx)
	if err := validateCtx.BindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transactionCategory := entities.TransactionCategory{
		Name: form.Name,
	}

	if err := c.service.CreateTransactionCategory(&transactionCategory, nil); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create TransactionCategory",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "TransactionCategory created successfully",
		"data":    (&model.TransactionCategory{}).Response(&transactionCategory),
	})
}

func (c *transactionCategoryController) Show(ctx *gin.Context) {
	id := ctx.Param("id")

	var transactionCategory entities.TransactionCategory

	err := c.service.FirstWithRelations(&transactionCategory, nil)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("TransactionCategory with ID %s not found", id),
		})
		return
	}

	ctx.JSON(http.StatusOK, (&model.TransactionCategory{}).Response(&transactionCategory))
}

// ต้องส่งมาทุก field
func (c *transactionCategoryController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	uintID, err := strconv.ParseUint(id, 10, 32) // แปลงเป็น uint32 ซึ่งสามารถแปลงเป็น uint ได้
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	form := model.TransactionCategoryQuery{}
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	transactionCategory := entities.TransactionCategory{
		Model: gorm.Model{ID: uint(uintID)},
		Name:  form.Name,
	}

	if err := c.service.UpdateTransactionCategory(&transactionCategory, nil); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update transactionCategory",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "transactionCategory updated successfully",
		"data":    (&model.TransactionCategory{}).Response(&transactionCategory),
	})
}

func (c *transactionCategoryController) Destroy(ctx *gin.Context) {
	id := ctx.Param("id")

	uintID, err := strconv.ParseUint(id, 10, 32) // แปลงเป็น uint32 ซึ่งสามารถแปลงเป็น uint ได้
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}
	transactionCategory := entities.TransactionCategory{
		Model: gorm.Model{ID: uint(uintID)},
	}

	if c.service.DeleteTransactionCategory(&transactionCategory); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete transactionCategory",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "transactionCategory deleted successfully",
	})
}
