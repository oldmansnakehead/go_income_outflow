package entities

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type CreditCard struct {
	gorm.Model

	Name        string          `gorm:"size:100;not null" json:"name"`
	CreditLimit decimal.Decimal `gorm:"type:decimal(20,2);default:0" json:"credit_limit"` // วงเงินของบัตรเครดิต
	Balance     decimal.Decimal `gorm:"type:decimal(20,2);default:0" json:"balance"`      // ยอดหนี้ที่ค้างจ่ายจากการใช้บัตร
	DueDate     uint            `gorm:"not null" json:"due_date"`                         // วันที่ครบกำหนดชำระ

	UserID uint `gorm:"not null" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user"`

	CreditCardDebts []CreditCardDebt `gorm:"foreignKey:CreditCardID" json:"credit_card_debts"`
	Transactions    []Transaction    `gorm:"polymorphic:Transactionable;" json:"transactions"`
}
