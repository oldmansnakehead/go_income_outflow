package model

import (
	"time"

	"go_income_outflow/pkg/model/common"
)

type CreditCardDebt struct {
	common.Model

	Amount           float64   // ยอดหนี้ที่ค้างชำระ
	DueDate          time.Time // วันที่ครบกำหนดชำระ
	PaidAmount       float64   // จำนวนเงินที่ชำระไปแล้ว
	RemainingBalance float64   // ยอดหนี้ที่เหลือ

	CreditCardID uint
	CreditCard   CreditCard
}
