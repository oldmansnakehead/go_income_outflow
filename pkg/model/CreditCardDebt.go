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

	PaymentTransactionID uint
	PaymentTransaction   Transaction
}

type CreditCardDebtRequest struct {
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`

	UserID           uint   `json:"user_id"`
	CreditCardID     uint   `json:"credit_card_id"`
	Date             string `json:"date"`
	InstallmentCount uint   `json:"installment_count"`
}
