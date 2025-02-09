package controller

// controller = handler

import (
	"fmt"
	"go_income_outflow/entities"
	"go_income_outflow/pkg/controller/common"
	"go_income_outflow/pkg/model"
	"go_income_outflow/pkg/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type (
	CreditCardControllerUseCase interface {
		common.ControllerUseCase
	}

	creditCardController struct {
		service service.CreditCardServiceUseCase
	}
)

func NewCreditCardController(service service.CreditCardServiceUseCase) CreditCardControllerUseCase {
	return &creditCardController{service: service}
}

func (c *creditCardController) Index(ctx *gin.Context) {
}

func (c *creditCardController) Store(ctx *gin.Context) {
	var form model.CreditCardRequest
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var creditCard entities.CreditCard
	if err := copier.Copy(&creditCard, &form); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to map form to credit card",
		})
		return
	}

	if err := c.service.CreateCreditCard(&creditCard, form.With); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create credit card",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Credit card created successfully",
		"data":    (&model.CreditCard{}).Response(&creditCard),
	})
}

func (c *creditCardController) Show(ctx *gin.Context) {
	id := ctx.Param("id")

	var form struct {
		With []string `json:"with" query:"with"`
	}

	if err := ctx.ShouldBindJSON(&form); err != nil {
		withArray := ctx.QueryArray("with[]")
		if len(withArray) > 0 {
			form.With = withArray
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid relations format",
			})
			return
		}
	}

	var creditCard entities.CreditCard

	err := c.service.FirstWithRelations(&creditCard, form.With)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Account with ID %s not found", id),
		})
		return
	}

	ctx.JSON(http.StatusOK, (&model.CreditCard{}).Response(&creditCard))
}

func (c *creditCardController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	uintID, err := strconv.ParseUint(id, 10, 32) // แปลงเป็น uint32 ซึ่งสามารถแปลงเป็น uint ได้
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	form := model.CreditCardRequest{}
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	creditCard := entities.CreditCard{
		Model: gorm.Model{ID: uint(uintID)}, // กำหนด ID โดยตรง เพราะมีการแปลงข้อมูล
	}
	if err := copier.Copy(&creditCard, &form); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to copy data into creditCard",
		})
		return
	}

	if err := c.service.UpdateCreditCard(&creditCard, form.With); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update credit card",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "credit card updated successfully",
		"data":    (&model.CreditCard{}).Response(&creditCard),
	})
}

func (c *creditCardController) Destroy(ctx *gin.Context) {
	id := ctx.Param("id")

	uintID, err := strconv.ParseUint(id, 10, 32) // แปลงเป็น uint32 ซึ่งสามารถแปลงเป็น uint ได้
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}
	creditCard := entities.CreditCard{
		Model: gorm.Model{ID: uint(uintID)},
	}

	if c.service.DeleteCreditCard(&creditCard); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete credit card",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "credit card deleted successfully",
	})
}
