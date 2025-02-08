package controller

// controller = handler

import (
	"fmt"
	"go_income_outflow/entities"
	"go_income_outflow/pkg/model"
	"go_income_outflow/pkg/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	if err := ac.service.CreateAccount(&account, form.With); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create account",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Account created successfully",
		"data":    (&model.Account{}).Response(&account), // keyword = struct literal
	})
}

func (ac *AccountController) Show(ctx *gin.Context) {
	id := ctx.Param("id")

	var form struct {
		With []string `json:"with" query:"with"` // รับจากทั้ง body และ query
	}

	withArray := ctx.QueryArray("with[]")
	if len(withArray) > 0 {
		form.With = withArray
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid relations format",
		})
		return
	}

	// Bind ข้อมูลจาก JSON body (ถ้ามี)
	if err := ctx.ShouldBindJSON(&form); err != nil {
		// ถ้าผิดพลาดในการ bind ให้รับข้อมูลจาก query parameter แทน
		if err := ctx.ShouldBindQuery(&form); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid relations format",
			})
			return
		}
	}

	var account entities.Account

	err := ac.service.GetAccountWithRelations(&account, form.With)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Account with ID %s not found", id),
		})
		return
	}

	ctx.JSON(http.StatusOK, (&model.Account{}).Response(&account))
}

// ต้องส่งมาทุก field
func (ac *AccountController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	uintID, err := strconv.ParseUint(id, 10, 32) // แปลงเป็น uint32 ซึ่งสามารถแปลงเป็น uint ได้
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	var form struct {
		Name   string   `json:"name"`
		UserID uint     `json:"user_id"`
		With   []string `json:"with"`
	}
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	account := entities.Account{
		Model:  gorm.Model{ID: uint(uintID)},
		Name:   form.Name,
		UserID: form.UserID,
	}

	if err := ac.service.UpdateAccount(&account, form.With); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update account",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Account updated successfully",
		"data":    (&model.Account{}).Response(&account),
	})
}

func (ac *AccountController) Destroy(ctx *gin.Context) {
	id := ctx.Param("id")

	uintID, err := strconv.ParseUint(id, 10, 32) // แปลงเป็น uint32 ซึ่งสามารถแปลงเป็น uint ได้
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}
	account := entities.Account{
		Model: gorm.Model{ID: uint(uintID)},
	}

	if ac.service.DeleteAccount(&account); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete account",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Account deleted successfully",
	})
}
