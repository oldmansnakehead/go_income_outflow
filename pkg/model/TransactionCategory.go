package model

import (
	"go_income_outflow/pkg/model/common"
)

type TransactionCategory struct {
	common.Model

	Name string
}
