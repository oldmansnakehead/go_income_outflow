// โครงสร้างข้อมูล layer ชั้นในสุด

package entities

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Name             string          `gorm:"size:50;not null" json:"name"`
	Amount           decimal.Decimal `gorm:"type:decimal(20,2);default:0" json:"amount"` // หรือเก็บเป็น string ก็ได้ใช้ shopstring/decimal คำนวนได้
	ExcludeFromTotal bool            `gorm:"default:false" json:"exclude_from_total"`
	Currency         string          `gorm:"size:3;not null" json:"currency"` // รหัสสกุลเงิน (เช่น "USD", "THB")

	UserID       uint          `json:"user_id"`
	User         User          `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"` // OnUpdate:CASCADE = user_id มีการเปลี่ยนแปลง account.user_id เปลี่ยนแปลงตาม // OnDelete CASCADE เหมือนกัน parent ลบ child โดนลบ
	Transactions []Transaction `gorm:"polymorphic:Transactionable;" json:"transactions"`
}
