package model

import (
	"go_income_outflow/entities"
	"go_income_outflow/pkg/model/common"
	"log"
	"time"

	"github.com/jinzhu/copier"
	"github.com/shopspring/decimal"
)

type (
	Account struct {
		common.Model

		Name             string
		Balance          decimal.Decimal
		ExcludeFromTotal bool
		Currency         string

		UserID uint
		User   User `gorm:"foreignKey:UserID"`

		Transactions []Transaction `gorm:"polymorphic:Transactionable;"`
	}

	AccountRequest struct {
		Name             string          `json:"name" binding:"required"`
		Balance          decimal.Decimal `json:"balance"`
		ExcludeFromTotal bool            `json:"exclude_from_total"`
		Currency         string          `json:"currency"`
		UserID           uint            `json:"user_id"`
		With             []string        `json:"with"`
	}

	AccountResponse struct {
		ID        uint      `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`

		Name             string          `json:"name"`
		Balance          decimal.Decimal `json:"balance"`
		ExcludeFromTotal bool            `json:"exclude_from_total"`
		Currency         string          `json:"currency"`

		UserID uint         `json:"user_id"`
		User   UserResponse `json:"user"`

		Transactions []TransactionResponse `json:"transactions"`
	}

	AccountQuery struct {
		Name   string   `json:"name"`
		UserID uint     `json:"user_id"`
		With   []string `json:"with"`
	}
)

func (r *Account) ToResponse(account *entities.Account) AccountResponse {
	var response AccountResponse
	err := copier.Copy(&response, account)
	if err != nil {
		log.Println("Error copying data:", err)
	}

	err = copier.Copy(&response.User, &account.User)
	if err != nil {
		log.Println("Error copying user data:", err)
	}

	err = copier.Copy(&response.Transactions, &account.Transactions)
	if err != nil {
		log.Println("Error copying user data:", err)
	}

	return response
}

func (r *Account) Response(account *entities.Account) AccountResponse {
	return r.ToResponse(account)
}
