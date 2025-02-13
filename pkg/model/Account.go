package model

import (
	"go_income_outflow/entities"
	"go_income_outflow/pkg/model/common"
	"time"
)

type (
	Account struct {
		common.Model

		Name string

		UserID *uint
		User   User `gorm:"foreignKey:UserID"`
	}

	AccountRequest struct {
		Name   string   `json:"name" binding:"required"`
		UserID uint     `json:"user_id"`
		With   []string `json:"with"`
	}

	AccountResponse struct {
		ID        uint      `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`

		Name   string       `json:"name"`
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
	r.UserID = account.UserID
	var userID uint
	if account.UserID != nil {
		userID = *account.UserID
	}

	r.User = User{
		ID:    userID,
		Name:  account.User.Name,
		Email: account.User.Email,
	}

	return r
}

func (r *Account) ToResponse() AccountResponse {
	var userID uint
	if r.UserID != nil {
		userID = *r.UserID
	}

	return AccountResponse{
		ID:     r.ID,
		Name:   r.Name,
		UserID: userID,
		User: UserResponse{
			ID:    userID,
			Name:  r.User.Name,
			Email: r.User.Email,
		},
	}
}

func (r *Account) Response(account *entities.Account) AccountResponse {
	return r.EntitiesToModel(account).ToResponse()
}
