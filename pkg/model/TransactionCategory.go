package model

import (
	"go_income_outflow/entities"
	"go_income_outflow/pkg/model/common"
	"time"
)

type (
	TransactionCategory struct {
		common.Model

		Name string
	}

	TransactionCategoryRequest struct {
		Name string `json:"name" binding:"required"`
	}

	TransactionCategoryResponse struct {
		ID        uint      `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`

		Name string `json:"name"`
	}

	TransactionCategoryQuery struct {
		Name string `json:"name"`
	}
)

func (r *TransactionCategory) EntitiesToModel(TransactionCategory *entities.TransactionCategory) *TransactionCategory {
	r.ID = TransactionCategory.ID
	r.CreatedAt = TransactionCategory.CreatedAt
	r.UpdatedAt = TransactionCategory.UpdatedAt
	r.Name = TransactionCategory.Name

	return r
}

func (r *TransactionCategory) ToResponse() TransactionCategoryResponse {
	return TransactionCategoryResponse{
		ID:   r.ID,
		Name: r.Name,
	}
}

func (r *TransactionCategory) Response(transactionCategory *entities.TransactionCategory) TransactionCategoryResponse {
	return r.EntitiesToModel(transactionCategory).ToResponse()
}
