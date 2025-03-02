package model

import (
	"time"

	"go_income_outflow/pkg/model/common"

	"github.com/shopspring/decimal"
)

type CreditCardDebt struct {
	common.Model

	Description string
	Amount      decimal.Decimal
	DueDate     time.Time

	CreditCardID uint
	CreditCard   CreditCard

	TransactionID uint
	Transaction   Transaction
}
