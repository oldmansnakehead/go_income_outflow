package controller

import (
	"go_income_outflow/db"
	"go_income_outflow/pkg/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct{}

func (u User) Create(ctx *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if ctx.Bind(&body) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user := model.User{Email: body.Email, Password: string(hash)}
	result := db.Conn.Create(&user)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"notice": "Success",
	})
}
