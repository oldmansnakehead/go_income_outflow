package request

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type (
	CustomRequest interface {
		BindJSON(obj any) error
		Bind(obj any) error
	}

	customRequest struct {
		ctx       *gin.Context
		validator *validator.Validate
	}
)

var (
	once              sync.Once
	validatorInstance *validator.Validate
)

func NewCustomRequest(req *gin.Context) CustomRequest {
	once.Do(func() {
		validatorInstance = validator.New()

		validatorInstance.RegisterValidation("uniqueEmail", uniqueUserEmailValidator)
	})

	return &customRequest{
		ctx:       req,
		validator: validatorInstance,
	}
}

func (r *customRequest) BindJSON(obj any) error {
	if err := r.ctx.ShouldBindJSON(&obj); err != nil {
		return err
	}
	fmt.Println(obj, 123456)
	if err := r.validator.Struct(obj); err != nil {
		return err
	}

	return nil
}

func (r *customRequest) Bind(obj any) error {
	if err := r.ctx.Bind(&obj); err != nil {
		return err
	}

	if err := r.validator.Struct(obj); err != nil {
		return err
	}

	return nil
}
