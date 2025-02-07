package model

import (
	"go_income_outflow/constants"
	"go_income_outflow/pkg/model/common"
	"time"
)

type Transaction struct {
	common.Model

	Amount          float64                   // จำนวนเงิน
	Type            constants.TransactionType // ประเภทของธุรกรรม (รายรับหรือรายจ่าย)
	TransactionDate time.Time                 // วันที่เกิดธุรกรรม
	Description     string                    // รายละเอียดเพิ่มเติม

	UserID uint
	User   User

	CategoryID uint
	Category   TransactionCategory

	AccountID uint
	Account   Account
}
