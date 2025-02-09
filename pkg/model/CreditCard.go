package model

import (
	"go_income_outflow/entities"
	"go_income_outflow/pkg/model/common"
	"time"
)

type (
	CreditCard struct {
		common.Model

		Name        string
		CreditLimit float64   // วงเงินของบัตรเครดิต
		Balance     float64   // ยอดหนี้ที่ค้างจ่ายจากการใช้บัตร
		DueDate     time.Time // วันที่ครบกำหนดชำระ

		UserID uint
		User   User `gorm:"foreignKey:UserID"`
	}

	CreditCardRequest struct {
		Name        string    `json:"name" binding:"required"`
		CreditLimit float64   `json:"credit_limit" binding:"required"`
		Balance     float64   `json:"balance"`
		DueDate     time.Time `json:"due_date"`

		UserID uint     `json:"user_id" binding:"required"`
		With   []string `json:"with"`
	}

	CreditCardResponse struct {
		ID          uint      `json:"id"`
		Name        string    `json:"name"`
		CreditLimit float64   `json:"credit_limit"`
		Balance     float64   `json:"balance"`
		DueDate     time.Time `json:"due_date"`

		UserID uint         `json:"user_id"`
		User   UserResponse `json:"user"`
	}
)

func (r *CreditCard) EntitiesToModel(creditCard *entities.CreditCard) *CreditCard {
	r.ID = creditCard.ID
	r.CreatedAt = creditCard.CreatedAt
	r.UpdatedAt = creditCard.UpdatedAt
	r.Name = creditCard.Name
	r.UserID = creditCard.UserID
	r.User = User{
		ID:    creditCard.UserID,
		Name:  creditCard.User.Name,
		Email: creditCard.User.Email,
	}

	return r
}

func (r *CreditCard) ToResponse() CreditCardResponse {
	return CreditCardResponse{
		ID:     r.ID,
		Name:   r.Name,
		UserID: r.UserID,
		User: UserResponse{
			ID:    r.UserID,
			Name:  r.User.Name,
			Email: r.User.Email,
		},
	}
}

func (r *CreditCard) Response(creditCard *entities.CreditCard) CreditCardResponse {
	return r.EntitiesToModel(creditCard).ToResponse()
}
