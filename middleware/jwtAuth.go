package middleware

import (
	"fmt"
	"go_income_outflow/db"
	"go_income_outflow/entities"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(ctx *gin.Context) {
	tokenString, err := ctx.Cookie("Authorization")
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		var user entities.User
		db.Conn.First(&user, claims["sub"])
		if user.ID == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		// เก็บ user ลงใน context
		ctx.Set("user", user)

		ctx.Next()
	} else {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
}
