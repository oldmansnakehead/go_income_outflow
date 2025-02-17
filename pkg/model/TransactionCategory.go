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
		Type bool
	}

	TransactionCategoryRequest struct {
		Name string `json:"name" binding:"required"`
		Type *bool  `json:"type" binding:"required"` // binding:"required" ใช้กับ bool ไม่ได้ต้อง *bool
	}

	TransactionCategoryResponse struct {
		ID        uint      `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`

		Name string `json:"name"`
		Type bool   `json:"type"`
	}

	TransactionCategoryQuery struct {
		Name string `json:"name"`
		Type bool   `json:"type"`
	}
)

func (r *TransactionCategory) EntitiesToModel(transactionCategory *entities.TransactionCategory) *TransactionCategory {
	r.ID = transactionCategory.ID
	r.CreatedAt = transactionCategory.CreatedAt
	r.UpdatedAt = transactionCategory.UpdatedAt
	r.Name = transactionCategory.Name
	r.Type = transactionCategory.Type

	return r
}

func (r *TransactionCategory) ToResponse() TransactionCategoryResponse {
	return TransactionCategoryResponse{
		ID:   r.ID,
		Name: r.Name,
		Type: r.Type,
	}
}

func (r *TransactionCategory) Response(transactionCategory *entities.TransactionCategory) TransactionCategoryResponse {
	return r.EntitiesToModel(transactionCategory).ToResponse()
}
