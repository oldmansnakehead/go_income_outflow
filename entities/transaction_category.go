package entities

import "gorm.io/gorm"

type TransactionCategory struct {
	gorm.Model
	Name string `gorm:"size:100;not null;"`
	Type bool   `gorm:"not null;default:false"` // false = รายจ่าย , true = รายรับ
}
