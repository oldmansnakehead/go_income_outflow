package entities

import (
	"time"
)

type RefreshToken struct {
	ID        string    `gorm:"type:varchar(36);primaryKey"`
	Token     string    `gorm:"type:text;not null"`
	UserID    uint      `gorm:"not null"`
	ExpiresAt time.Time `gorm:"type:timestamp;not null"`
	Counter   uint      `gorm:"default:0"` // counter สำหรับนับจำนวนการใช้ refresh token
	Revoke    bool      `gorm:"default:false"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
