// โครงสร้างข้อมูล layer ชั้นในสุด

package entities

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Name string `gorm:"size:50;not null"; json:name`

	UserID uint
	User   User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // OnUpdate:CASCADE = user_id มีการเปลี่ยนแปลง account.user_id เปลี่ยนแปลงตาม // OnDelete CASCADE เหมือนกัน parent ลบ child โดนลบ
}
