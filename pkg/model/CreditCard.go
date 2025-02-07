package model

import (
	"go_income_outflow/pkg/model/common"
	"time"
)

type CreditCard struct {
	common.Model

	CardName    string    `gorm:"size:100;not null"`
	CreditLimit float64   `gorm:"not null"`           // วงเงินของบัตรเครดิต
	Balance     float64   `gorm:"default:0;not null"` // ยอดหนี้ที่ค้างจ่ายจากการใช้บัตร
	DueDate     time.Time `gorm:"not null"`           // วันที่ครบกำหนดชำระ

	UserID uint `gorm:"not null"`
	User   User
}
