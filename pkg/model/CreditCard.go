package model

import (
	"go_income_outflow/entities"
	"go_income_outflow/pkg/model/common"
	"log"

	"github.com/jinzhu/copier"
	"github.com/shopspring/decimal"
)

type (
	CreditCard struct {
		common.Model

		Name        string
		CreditLimit decimal.Decimal // วงเงินของบัตรเครดิต
		Balance     decimal.Decimal // ยอดหนี้ที่ค้างจ่ายจากการใช้บัตร
		DueDate     uint            // วันที่ครบกำหนดชำระ

		UserID uint
		User   User `gorm:"foreignKey:UserID"`

		Transactions []Transaction `gorm:"polymorphic:Transactionable;"`
	}

	CreditCardRequest struct {
		Name        string          `json:"name" binding:"required"`
		CreditLimit decimal.Decimal `json:"credit_limit" binding:"required"`
		Balance     decimal.Decimal `json:"balance"`
		DueDate     uint            `json:"due_date"`

		UserID uint     `json:"user_id" binding:"required"`
		With   []string `json:"with"`
	}

	CreditCardResponse struct {
		ID          uint            `json:"id"`
		Name        string          `json:"name"`
		CreditLimit decimal.Decimal `json:"credit_limit"`
		Balance     decimal.Decimal `json:"balance"`
		DueDate     uint            `json:"due_date"`

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
