package model

import (
	"go_income_outflow/pkg/model/common"
)

type User struct {
	common.Model

	Email    string
	Password string
	Name     string
}
