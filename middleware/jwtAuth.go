package middleware

import (
	"go_income_outflow/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth(ctx *gin.Context) {
	ac, err := token.ExtractAccessClaims(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	/* if _, err := Authorize(&ac.User, c.Request.Method, c.Request.URL.Path); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized to access: " + err.Error(),
		})
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	} */

	ctx.Set("user", &ac.User)
	ctx.Next()
}
