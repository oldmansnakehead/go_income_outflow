package token

import (
	"fmt"
	"go_income_outflow/entities"
	"go_income_outflow/helpers"
	"go_income_outflow/pkg/model"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RefreshClaims struct {
	TokenID string `json:"jti"`
	Counter uint
	User    model.User
	jwt.RegisteredClaims
}

func NewRefreshClaims(u *model.User) (rc *RefreshClaims) {
	rc = &RefreshClaims{User: *u, TokenID: uuid.New().String()}
	n := now()
	rc.IssuedAt = &jwt.NumericDate{Time: n} //
	rc.ExpiresAt = &jwt.NumericDate{Time: n.Add(refreshExp)}
	rc.Counter = 0
	rc.UpdateTime()

	return rc
}

func (rc *RefreshClaims) UpdateTime() *RefreshClaims {
	now := now()
	rc.ExpiresAt = &jwt.NumericDate{Time: now.Add(refreshExp)}
	rc.NotBefore = &jwt.NumericDate{Time: now.Add(refreshNBF)}
	return rc
}

func (rc *RefreshClaims) JwtString() (refreshTokenString string, err error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rc) // สร้าง token และ บันทึก claims เข้า payload
	refreshTokenString, err = refreshToken.SignedString([]byte(os.Getenv("REFRESH_SECRET_KEY")))
	if err != nil {
		return "", fmt.Errorf("error refresh token SignedString")
	}

	return refreshTokenString, nil
}

func ExtractRefreshClaims(c *gin.Context, option ...jwt.ParserOption) (*RefreshClaims, error) {
	tokenString := helpers.ExtractRefreshToken(c)

	// not before เช็คตรงนี้ด้วย
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET_KEY")), nil
	}, option...)

	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %v", err)
	}

	rc, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token: %v", "invlid claims")
	}

	if rc.IsExpired() {
		return nil, fmt.Errorf("invalid refresh token: %v", "expired")
	}

	return rc, nil
}

func (rc *RefreshClaims) IsExpired() bool {
	now := now()
	exp := rc.ExpiresAt.Time
	return now.After(exp)
}

func (rc *RefreshClaims) Rotate(r *entities.RefreshToken) (err error) {
	now := now()

	// check expire
	if !now.Before(rc.ExpiresAt.Time) {
		return fmt.Errorf("token expired: %s", rc.ExpiresAt.Time.String())
	}

	// check not before
	if now.Before(rc.NotBefore.Time) {
		return fmt.Errorf("token not before: %s", rc.NotBefore.Time.String())
	}

	// check revoke
	if r.Revoke {
		return fmt.Errorf("invalid token revoke: %s", r.ID)
	}

	// check counter
	if r.Counter > rc.Counter {
		return ErrInvalidTokenCounter
	}

	rc.Counter++
	rc.UpdateTime()

	return nil
}
