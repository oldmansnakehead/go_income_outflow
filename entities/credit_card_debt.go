package entities

import (
	"time"

	"gorm.io/gorm"
)

type CreditCardDebt struct {
	gorm.Model

	Amount           float64   `gorm:"not null"`  // ยอดหนี้ที่ค้างชำระ
	DueDate          time.Time `gorm:"not null"`  // วันที่ครบกำหนดชำระ
	PaidAmount       float64   `gorm:"default:0"` // จำนวนเงินที่ชำระไปแล้ว
	RemainingBalance float64   `gorm:"not null"`  // ยอดหนี้ที่เหลือ

	CreditCardID uint       `gorm:"not null"`
	CreditCard   CreditCard `gorm:"foreignKey:CreditCardID"`
}
