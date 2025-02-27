package model

import (
	"go_income_outflow/entities"
	"log"
	"time"

	"github.com/jinzhu/copier"
)

type (
	TransactionCategory struct {
		ID        uint
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt time.Time

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

func (r *TransactionCategory) ToResponse(category *entities.TransactionCategory) TransactionCategoryResponse {
	var response TransactionCategoryResponse
	err := copier.Copy(&response, category)
	if err != nil {
		log.Println("Error copying data:", err)
	}

	return response
}

func (r *TransactionCategory) Response(transactionCategory *entities.TransactionCategory) TransactionCategoryResponse {
	return r.ToResponse(transactionCategory)
}
