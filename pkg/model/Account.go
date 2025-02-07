package model

import "go_income_outflow/pkg/model/common"

type Account struct {
	common.Model

	Name string

	UserID uint
	User   User `gorm:"foreignKey:UserID"`
}

type AccountRequest struct {
	Name   string `json:"name" binding:"required"`
	UserID uint   `json:"user_id"`
}

type AccountResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
