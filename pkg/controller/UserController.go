package controller

import (
	"errors"
	"fmt"
	"go_income_outflow/pkg/controller/common"
	"go_income_outflow/pkg/custom/request"
	"go_income_outflow/pkg/model"
	"go_income_outflow/pkg/service"
	"go_income_outflow/token"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type (
	UserControllerUseCase interface {
		common.ControllerUseCase
		TestAuth(ctx *gin.Context)
		Login(ctx *gin.Context)
		Logout(ctx *gin.Context)
		RefreshToken(ctx *gin.Context)
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
	refreshToken, accessToken, refreshExpiresAt, accessExpiresAt, err := u.service.Login(itemReq.Email, itemReq.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// ส่ง response พร้อม token
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", accessToken, int(time.Until(accessExpiresAt).Seconds()), "", "", false, true)
	ctx.SetCookie("RefreshToken", refreshToken, int(time.Until(refreshExpiresAt).Seconds()), "", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"notice":        "Login successful",
		"refresh_token": refreshToken,
		"access_token":  accessToken,
	})
}

func (u *userController) Logout(ctx *gin.Context) {
	// ดึง claims เพื่อเอา user
	rc, err := token.ExtractRefreshClaims(ctx, jwt.WithoutClaimsValidation())
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := u.service.Logout(rc.User.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ลบหรือยกเลิก Cookie
	ctx.SetCookie("Authorization", "", -1, "", "", false, true)
	ctx.SetCookie("RefreshToken", "", -1, "", "", false, true)
	ctx.Set("user", nil)

	ctx.JSON(http.StatusOK, gin.H{
		"notice": "Logout successful",
	})
}

func (uc *userController) TestAuth(ctx *gin.Context) {
	user, _ := ctx.Get("user")

	ctx.JSON(http.StatusOK, gin.H{
		"notice": user,
	})
}

func (u *userController) RefreshToken(ctx *gin.Context) {
	rc, err := token.ExtractRefreshClaims(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	r, err := u.service.GetRefreshTokenByID(rc.TokenID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = rc.Rotate(r)
	if err != nil {
		if errors.Is(err, token.ErrInvalidTokenCounter) {
			// revoke
			if err := u.service.RevokeRefreshToken(rc.TokenID); err != nil {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"error": err.Error(),
				})
			} else {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid token counter: token has been revoked",
				})
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	refreshTokenString, err := rc.JwtString()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	user := &model.User{
		ID:    rc.User.ID,
		Name:  rc.User.Name,
		Email: rc.User.Email,
	}
	ac := token.NewAccessClaims(user)
	accessTokenString, err := ac.JwtString()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := u.service.UpdateRefreshToken(rc, refreshTokenString); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("failed to update refresh token: %v", err),
		})
		return
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", accessTokenString, int(time.Until(ac.ExpiresAt.Time).Seconds()), "", "", false, true)
	ctx.SetCookie("RefreshToken", refreshTokenString, int(time.Until(rc.ExpiresAt.Time).Seconds()), "", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"refresh_token": refreshTokenString,
		"access_token":  accessTokenString,
	})
}
