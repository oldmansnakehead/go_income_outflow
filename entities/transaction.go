package entities

import (
	"go_income_outflow/constants"
	"time"

	"gorm.io/gorm"
)

// ต้องการ Enum type transaction_type

type Transaction struct {
	gorm.Model

	Amount          float64                   `gorm:"not null"`                                      // จำนวนเงิน
	Type            constants.TransactionType `gorm:"type:transaction_type;default:'ANY';not null;"` // ประเภทของธุรกรรม (รายรับหรือรายจ่าย)
	TransactionDate time.Time                 `gorm:"not null"`                                      // วันที่เกิดธุรกรรม
	Description     string                    `gorm:"type:text"`                                     // รายละเอียดเพิ่มเติม

	UserID uint `gorm:"not null"`
	User   User

	CategoryID uint                `gorm:"not null"`
	Category   TransactionCategory `gorm:"foreignKey:CategoryID"`

	AccountID uint    `gorm:"not null"`
	Account   Account `gorm:"foreignKey:AccountID"`
}
