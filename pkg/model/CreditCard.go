package model

import (
	"go_income_outflow/entities"
	"go_income_outflow/pkg/model/common"
	"log"
	"time"

	"github.com/jinzhu/copier"
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

		Transactions []Transaction `gorm:"polymorphic:Transactionable;"`
	}

	CreditCardRequest struct {
		Name        string  `json:"name" binding:"required"`
		CreditLimit float64 `json:"credit_limit" binding:"required"`
		Balance     float64 `json:"balance"`
		DueDate     string  `json:"due_date"`

		UserID uint     `json:"user_id" binding:"required"`
		With   []string `json:"with"`
	}

	CreditCardResponse struct {
		ID          uint      `json:"id"`
		Name        string    `json:"name"`
		CreditLimit float64   `json:"credit_limit"`
		Balance     float64   `json:"balance"`
		DueDate     time.Time `json:"due_date"`

		UserID       uint         `json:"user_id"`
		User         UserResponse `json:"user"`
		Transactions []Transaction
	}
)

func (r *CreditCard) ToResponse(creditCardResponse *entities.CreditCard) CreditCardResponse {
	var response CreditCardResponse
	err := copier.Copy(&response, creditCardResponse)
	if err != nil {
		log.Println("Error copying data:", err)
	}

	err = copier.Copy(&response.User, &creditCardResponse.User)
	if err != nil {
		log.Println("Error copying user data:", err)
	}

	return response
}

func (r *CreditCard) Response(creditCard *entities.CreditCard) CreditCardResponse {
	return r.ToResponse(creditCard)
}
