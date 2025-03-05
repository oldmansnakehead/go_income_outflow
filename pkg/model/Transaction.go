package model

import (
	"go_income_outflow/entities"
	"log"
	"time"

	"github.com/jinzhu/copier"
	"github.com/shopspring/decimal"
)

type (
	Transaction struct {
		ID        uint
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt time.Time

		Date        string
		Amount      decimal.Decimal // จำนวนเงิน
		Description string          // รายละเอียดเพิ่มเติม

		UserID uint
		User   User

		CategoryID uint
		Category   TransactionCategory

		TransactionableID   uint
		TransactionableType string
		Transactionable     interface{} `gorm:"-"`
	}

	TransactionRequest struct {
		Date        string          `json:"date" binding:"required"`
		Amount      decimal.Decimal `json:"amount" binding:"required"`
		Description string          `json:"description"`

		UserID uint `json:"user_id"`
		User   User

		CategoryID uint `json:"category_id"`

		TransactionableID   uint   `gorm:"not null" json:"transactionable_id"`
		TransactionableType string `gorm:"not null" json:"transactionable_type"`

		CreditCardDebtId uint `json:"credit_card_debt_id"`

		InstallmentCount uint `json:"installment_count"`

		With []string `json:"with"`
	}

	TransactionResponse struct {
		ID        uint      `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`

		Date   string          `json:"date"`
		Amount decimal.Decimal `json:"amount"`

		Description string `json:"description"`

		UserID uint         `json:"user_id"`
		User   UserResponse `json:"user"`

		CategoryID uint                        `json:"category_id"`
		Category   TransactionCategoryResponse `json:"category"`

		TransactionableID   uint        `gorm:"not null" json:"transactionable_id"`
		TransactionableType string      `gorm:"not null" json:"transactionable_type"`
		Transactionable     interface{} `json:"transactionable"`
	}

	TransactionQuery struct {
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`

		Date string `json:"date"`

		UserID     uint `json:"user_id"`
		CategoryID uint `json:"category_id"`

		TransactionableID   uint   `gorm:"not null" json:"transactionable_id"`
		TransactionableType string `gorm:"not null" json:"transactionable_type"`

		With []string `json:"with"`
	}
)

func (r *Transaction) ToResponse(transaction *entities.Transaction) TransactionResponse {
	var response TransactionResponse
	err := copier.Copy(&response, transaction)
	if err != nil {
		log.Println("Error copying data:", err)
	}

	err = copier.Copy(&response.User, &transaction.User)
	if err != nil {
		log.Println("Error copying user data:", err)
	}

	err = copier.Copy(&response.Category, &transaction.Category)
	if err != nil {
		log.Println("Error copying user data:", err)
	}

	return response
}

func (r *Transaction) Response(transaction *entities.Transaction) TransactionResponse {
	return r.ToResponse(transaction)
}
