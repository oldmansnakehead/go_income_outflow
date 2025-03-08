package repository

import (
	"go_income_outflow/entities"
	"time"

	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	SaveRefreshToken(userID uint, token string, expiresAt time.Time, tokenID string, tx *gorm.DB) error
	ClearRefreshToken(userID uint, tx *gorm.DB) error
	IncrementCounter(tokenID string) error
}

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

// SaveRefreshToken บันทึก refresh token ลงในฐานข้อมูล
func (r *refreshTokenRepository) SaveRefreshToken(userID uint, token string, expiresAt time.Time, tokenID string, tx *gorm.DB) error {
	conn := r.db
	if tx != nil {
		conn = tx
	}

	refreshToken := &entities.RefreshToken{
		ID:        tokenID,
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
		Counter:   0, // เริ่มต้น counter เป็น 0
	}
	result := conn.Create(refreshToken)
	return result.Error
}

func (r *refreshTokenRepository) ClearRefreshToken(userID uint, tx *gorm.DB) error {
	conn := r.db
	if tx != nil {
		conn = tx
	}

	return conn.Model(&entities.RefreshToken{}).Where("user_id = ?", userID).Update("revoke", true).Error
}

// IncrementCounter เพิ่มค่า counter เมื่อมีการใช้ refresh token
func (r *refreshTokenRepository) IncrementCounter(tokenID string) error {
	return r.db.Model(&entities.RefreshToken{}).
		Where("token_id = ?", tokenID).
		Update("counter", gorm.Expr("counter + 1")).Error
}
