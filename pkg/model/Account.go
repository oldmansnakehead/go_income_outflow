package model

import (
	"go_income_outflow/pkg/model/common"
)

type Account struct {
	common.Model

	Name string

	UserID uint
	User   User
}
