package entities

import "gorm.io/gorm"

type TransactionCategory struct {
	gorm.Model
	Name string `gorm:"size:100;not null;" json:"name"`
	Type bool   `gorm:"not null;default:false" json:"type"` // false = รายจ่าย , true = รายรับ
}
