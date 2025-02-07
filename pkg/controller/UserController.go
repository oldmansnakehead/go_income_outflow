package controller

import (
	"go_income_outflow/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{service: service}
}

func (uc *UserController) Index(ctx *gin.Context) {
}

func (u *UserController) Store(ctx *gin.Context) {
	var body service.CreateUserBody

	if ctx.Bind(&body) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	err := u.service.CreateUser(&body)
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

func (uc *UserController) Show(ctx *gin.Context) {
}

func (uc *UserController) Update(ctx *gin.Context) {
}

func (uc *UserController) Destroy(ctx *gin.Context) {
}

func (u *UserController) Login(ctx *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	// รับข้อมูลจาก body
	if ctx.Bind(&body) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// เรียกใช้ service เพื่อทำการล็อกอิน
	token, err := u.service.Login(body.Email, body.Password)
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

func (uc *UserController) TestAuth(ctx *gin.Context) {
	user, _ := ctx.Get("user")

	ctx.JSON(http.StatusOK, gin.H{
		"notice": user,
	})
}
