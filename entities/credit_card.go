package entities

import (
	"time"

	"gorm.io/gorm"
)

type CreditCard struct {
	gorm.Model

	Name        string    `gorm:"size:100;not null" json:"name"`
	CreditLimit float64   `gorm:"not null" json:"credit_limit"`       // วงเงินของบัตรเครดิต
	Balance     float64   `gorm:"default:0;not null" json:"balance"`  // ยอดหนี้ที่ค้างจ่ายจากการใช้บัตร
	DueDate     time.Time `gorm:"type:date;not null" json:"due_date"` // วันที่ครบกำหนดชำระ

	UserID uint `gorm:"not null" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user"`

	Transactions []Transaction `gorm:"polymorphic:Transactionable;" json:"transactions"`
}
