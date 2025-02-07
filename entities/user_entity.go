package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string    `gorm:"type:varchar(128);unique;not null;"`
	Password string    `gorm:"size:255;not null;"`
	Name     string    `gorm:"size:50;not null;"`
	Accounts []Account `gorm:"foreignKey:UserID"`
}
