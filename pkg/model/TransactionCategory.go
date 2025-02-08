package model

import (
	"go_income_outflow/pkg/model/common"
)

type TransactionCategory struct {
	common.Model

	Name string
}

type TransactionCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type TransactionCategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
