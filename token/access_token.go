package token

import (
	"fmt"
	"go_income_outflow/helpers"
	"go_income_outflow/pkg/model"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	User model.User
	jwt.RegisteredClaims
}

func NewAccessClaims(u *model.User) (ac *AccessClaims) {
	ac = &AccessClaims{User: *u}
	n := now()
	ac.IssuedAt = jwt.NewNumericDate(n)
	ac.ExpiresAt = jwt.NewNumericDate(n.Add(accessExp))
	return
}

func (ac *AccessClaims) JwtString() (accessTokenString string, err error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, ac)
	accessTokenString, err = accessToken.SignedString([]byte(os.Getenv("ACCESS_SECRET_KEY")))
	if err != nil {
		return "", fmt.Errorf("error access token SignedString")
	}
	return accessTokenString, nil
}

func (ac *AccessClaims) IsExpired() bool {
	now := now()
	exp := ac.ExpiresAt.Time
	return now.After(exp)
}

func ExtractAccessClaims(ctx *gin.Context) (*AccessClaims, error) {
	tokenString := helpers.ExtractJWT(ctx)

	// แปลงเป็น token พร้อมดึง claims (ข้อมูลใน payload)
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET_KEY")), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid access token: %v", err)
	}

	// แปลง claims และตรวจสอบความถูกต้องของ Token
	ac, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid access token: %v", "invlid claims")
	}

	if ac.IsExpired() {
		return nil, fmt.Errorf("invalid access token: %v", "expired")
	}

	return ac, nil
}
