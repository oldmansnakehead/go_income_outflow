package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string    `gorm:"type:varchar(128);unique;not null;" json:"email"`
	Password string    `gorm:"size:255;not null;" json:"password"`
	Name     string    `gorm:"size:50;not null;" json:"name"`
	Accounts []Account `gorm:"foreignKey:UserID" json:"accounts"`
}
