package model

import (
	"go_income_outflow/entities"
	"go_income_outflow/pkg/model/common"
)

type (
	Account struct {
		common.Model

		Name string

		UserID uint
		User   User `gorm:"foreignKey:UserID"`
	}

	AccountRequest struct {
		Name   string   `json:"name" binding:"required"`
		UserID uint     `json:"user_id"`
		With   []string `json:"with"`
	}

	AccountResponse struct {
		ID     uint         `json:"id"`
		Name   string       `json:"name"`
		UserID uint         `json:"user_id"`
		User   UserResponse `json:"user"`
	}
)

func (ac *Account) EntitiesToModel(account *entities.Account) *Account {
	ac.ID = account.ID
	ac.CreatedAt = account.CreatedAt
	ac.UpdatedAt = account.UpdatedAt
	ac.Name = account.Name
	ac.UserID = account.UserID
	ac.User = User{
		ID:    account.UserID,
		Name:  account.User.Name,
		Email: account.User.Email,
	}

	return ac
}

func (ac *Account) ToResponse() AccountResponse {
	return AccountResponse{
		ID:     ac.ID,
		Name:   ac.Name,
		UserID: ac.UserID,
		User: UserResponse{
			ID:    ac.UserID,
			Name:  ac.User.Name,
			Email: ac.User.Email,
		},
	}
}

func (ac *Account) Response(account *entities.Account) AccountResponse {
	return ac.EntitiesToModel(account).ToResponse()
}
