package entities

import "gorm.io/gorm"

type TransactionCategory struct {
	gorm.Model
	Name string `gorm:"size:100;not null;"`
}
