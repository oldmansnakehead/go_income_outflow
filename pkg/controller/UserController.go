package controller

import (
	"go_income_outflow/pkg/controller/common"
	"go_income_outflow/pkg/custom/request"
	"go_income_outflow/pkg/model"
	"go_income_outflow/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	UserControllerUseCase interface {
		common.ControllerUseCase
		TestAuth(ctx *gin.Context)
		Login(ctx *gin.Context)
	}
	userController struct {
		service service.UserServiceUseCase
	}
)

func NewUserController(service service.UserServiceUseCase) UserControllerUseCase {
	return &userController{service: service}
}

func (uc *userController) Index(ctx *gin.Context) {
}

func (u *userController) Store(ctx *gin.Context) {
	var itemReq model.UserRequest

	validateCtx := request.NewCustomRequest(ctx)
	if err := validateCtx.BindJSON(&itemReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := u.service.CreateUser(&itemReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"notice": "Success",
	})
}

func (uc *userController) Show(ctx *gin.Context) {
}

func (uc *userController) Update(ctx *gin.Context) {
}

func (uc *userController) Destroy(ctx *gin.Context) {
}

func (u *userController) Login(ctx *gin.Context) {
	var itemReq struct {
		Email    string
		Password string
	}

	validateCtx := request.NewCustomRequest(ctx)
	if err := validateCtx.BindJSON(&itemReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// เรียกใช้ service เพื่อทำการล็อกอิน
	token, err := u.service.Login(itemReq.Email, itemReq.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// ส่ง response พร้อม token
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", token, 3600*24*30, "", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"notice": "Login successful",
	})
}

func (uc *userController) TestAuth(ctx *gin.Context) {
	user, _ := ctx.Get("user")

	ctx.JSON(http.StatusOK, gin.H{
		"notice": user,
	})
}
