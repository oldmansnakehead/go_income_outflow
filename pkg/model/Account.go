package model

import (
	"go_income_outflow/entities"
	"go_income_outflow/pkg/model/common"
	"time"

	"github.com/shopspring/decimal"
)

type (
	Account struct {
		common.Model

		Name             string
		Amount           decimal.Decimal
		ExcludeFromTotal bool
		Currency         string

		UserID uint
		User   User `gorm:"foreignKey:UserID"`
	}

	AccountRequest struct {
		Name             string          `json:"name" binding:"required"`
		Amount           decimal.Decimal `json:"amount"`
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
		Amount           decimal.Decimal `json:"amount"`
		ExcludeFromTotal bool            `json:"exclude_from_total"`
		Currency         string          `json:"currency"`

		UserID uint         `json:"user_id"`
		User   UserResponse `json:"user"`
	}

	AccountQuery struct {
		Name   string   `json:"name"`
		UserID uint     `json:"user_id"`
		With   []string `json:"with"`
	}
)

func (r *Account) EntitiesToModel(account *entities.Account) *Account {
	r.ID = account.ID
	r.CreatedAt = account.CreatedAt
	r.UpdatedAt = account.UpdatedAt
	r.Name = account.Name
	r.Amount = account.Amount
	r.ExcludeFromTotal = account.ExcludeFromTotal
	r.Currency = account.Currency
	r.UserID = account.UserID

	r.User = User{
		ID:    r.UserID,
		Name:  account.User.Name,
		Email: account.User.Email,
	}

	return r
}

func (r *Account) ToResponse() AccountResponse {
	return AccountResponse{
		ID:               r.ID,
		Name:             r.Name,
		UserID:           r.UserID,
		Amount:           r.Amount,
		ExcludeFromTotal: r.ExcludeFromTotal,
		Currency:         r.Currency,
		User: UserResponse{
			ID:    r.UserID,
			Name:  r.User.Name,
			Email: r.User.Email,
		},
	}
}

func (r *Account) Response(account *entities.Account) AccountResponse {
	return r.EntitiesToModel(account).ToResponse()
}
