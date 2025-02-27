package migrations

import (
	"go_income_outflow/db"
	"go_income_outflow/entities"
)

func Migrate() {
	// สร้าง type สำหรับ ENUM
	enumNames := []string{
		// รายชื่อ ENUM ที่ต้องการสร้าง
		"transaction_type",
	}
	initialEnum(enumNames)

	tx := db.Conn.Begin()

	if err := tx.AutoMigrate(
		&entities.User{},
		&entities.Account{},
		&entities.CreditCardDebt{},
		&entities.CreditCard{},
		&entities.TransactionCategory{},
		&entities.Transaction{},
	); err != nil {
		tx.Rollback()
		panic(err)
	}

	tx.Commit()
}
