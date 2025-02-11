package request

import (
	"go_income_outflow/db"
	"go_income_outflow/pkg/model"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func uniqueUserEmailValidator(fl validator.FieldLevel) bool {
	email := fl.Field().String()

	var user model.User
	if err := db.Conn.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return true
		}
	}
	// ถ้ามีอีเมลนี้ในฐานข้อมูลแล้ว
	return false
}
