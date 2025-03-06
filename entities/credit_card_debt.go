package entities

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type CreditCardDebt struct {
	gorm.Model

	Description string          `gorm:"size:100;not null" json:"name"`
	Amount      decimal.Decimal `gorm:"type:decimal(20,2);default:0" json:"amount"` // ยอดหนี้ที่ค้างชำระ
	DueDate     string          `gorm:"type:date;not null" json:"due_date"`

	CreditCardID uint       `gorm:"not null" json:"credit_card_id"`
	CreditCard   CreditCard `gorm:"foreignKey:CreditCardID" json:"credit_card"`

	TransactionID uint        `json:"transaction_id"`
	Transaction   Transaction `gorm:"foreignKey:TransactionID" json:"transaction"`

	PaymentTransactionID *uint       `json:"payment_transaction_id"`
	PaymentTransaction   Transaction `gorm:"foreignKey:PaymentTransactionID" json:"payment_transaction"`
}
