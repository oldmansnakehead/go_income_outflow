package migrations

import (
	"go_income_outflow/db"
	"go_income_outflow/entities"
)

func Migrate() {
	db.Conn.AutoMigrate(
		&entities.User{},
		&entities.Account{},
		&entities.CreditCardDebt{},
		&entities.CreditCard{},
		&entities.TransactionCategory{},
		&entities.Transaction{},
	)
}
