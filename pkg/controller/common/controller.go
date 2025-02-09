package common

import "github.com/gin-gonic/gin"

type ControllerUseCase interface {
	Index(ctx *gin.Context)
	Store(ctx *gin.Context)
	Show(ctx *gin.Context)
	Update(ctx *gin.Context)
	Destroy(ctx *gin.Context)
}
