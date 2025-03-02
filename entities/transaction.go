package entities

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ต้องการ Enum type transaction_type

type Transaction struct {
	gorm.Model

	Date        string          `gorm:"type:date" json:"date"`
	Amount      decimal.Decimal `gorm:"type:decimal(20,2);default:0" json:"amount"` // หรือเก็บเป็น string ก็ได้ใช้ shopstring/decimal คำนวนได้
	Description string          `gorm:"type:text" json:"description"`

	UserID uint `gorm:"not null" json:"user_id"`
	User   User `json:"user"`

	CategoryID uint                `gorm:"not null" json:"category_id"`
	Category   TransactionCategory `gorm:"foreignKey:CategoryID;references:ID" json:"category"`

	TransactionableID   uint        `gorm:"not null" json:"transactionable_id"`
	TransactionableType string      `gorm:"size:255;not null" json:"transactionable_type"` // "accounts" or "credit_cards"
	Transactionable     interface{} `gorm:"-" json:"transactionable"`
}
