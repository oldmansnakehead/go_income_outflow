package entities

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Name string `gorm:"size:50;not null;"`

	UserID uint
	User   User
}
