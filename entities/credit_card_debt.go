package entities

import (
	"time"

	"gorm.io/gorm"
)

type CreditCardDebt struct {
	gorm.Model

	Amount           float64   `gorm:"not null" json:"amount"`            // ยอดหนี้ที่ค้างชำระ
	DueDate          time.Time `gorm:"not null" json:"due_date"`          // วันที่ครบกำหนดชำระ
	PaidAmount       float64   `gorm:"default:0" json:"paid_amount"`      // จำนวนเงินที่ชำระไปแล้ว
	RemainingBalance float64   `gorm:"not null" json:"remaining_balance"` // ยอดหนี้ที่เหลือ

	CreditCardID uint       `gorm:"not null" json:"credit_card_id"`
	CreditCard   CreditCard `gorm:"foreignKey:CreditCardID" json:"credit_card"`
}
