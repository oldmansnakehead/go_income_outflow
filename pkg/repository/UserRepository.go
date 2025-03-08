package repository

import (
	"errors"
	"fmt"
	"go_income_outflow/entities"
	"go_income_outflow/token"
	"time"

	"gorm.io/gorm"
)

type (
	UserRepository interface {
		BeginTransaction() *gorm.DB
		CommitTransaction(tx *gorm.DB) error
		RollbackTransaction(tx *gorm.DB) error

		CreateUser(user *entities.User) error
		FindUserByEmail(email string) (*entities.User, error)

		SaveRefreshToken(userID uint, token string, expiresAt time.Time, tokenID string, tx *gorm.DB) error
		ClearRefreshToken(userID uint, tx *gorm.DB) error
		GetRefreshTokenByID(tokenID string) (*entities.RefreshToken, error)
		UpdateRefreshToken(rc *token.RefreshClaims, refreshToken string) error
		RevokeRefreshToken(tokenID string) error
	}

	userRepository struct {
		db               *gorm.DB
		refreshTokenRepo RefreshTokenRepository
	}
)

func NewUserRepository(db *gorm.DB, refreshTokenRepo RefreshTokenRepository) UserRepository {
	return &userRepository{db: db, refreshTokenRepo: refreshTokenRepo}
}

func (r *userRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *userRepository) CommitTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (r *userRepository) RollbackTransaction(tx *gorm.DB) error {
	return tx.Rollback().Error
}

func (r *userRepository) CreateUser(User *entities.User) error {
	return r.db.Create(User).Error
}

func (r *userRepository) FindUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	result := r.db.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) SaveRefreshToken(userID uint, token string, expiresAt time.Time, tokenID string, tx *gorm.DB) error {
	return r.refreshTokenRepo.SaveRefreshToken(userID, token, expiresAt, tokenID, tx)
}

func (r *userRepository) ClearRefreshToken(userID uint, tx *gorm.DB) error {
	return r.refreshTokenRepo.ClearRefreshToken(userID, tx)
}

func (r *userRepository) GetRefreshTokenByID(tokenID string) (*entities.RefreshToken, error) {
	var refreshToken entities.RefreshToken

	err := r.db.Where("id = ?", tokenID).First(&refreshToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, fmt.Errorf("failed to get refresh token: %v", err)
	}

	return &refreshToken, nil
}

func (r *userRepository) UpdateRefreshToken(rc *token.RefreshClaims, refreshToken string) error {
	// อัปเดทค่า counter และ expires_at
	result := r.db.Model(&entities.RefreshToken{}).
		Where("id = ?", rc.TokenID).
		Updates(map[string]interface{}{
			"token":      refreshToken,
			"expires_at": rc.ExpiresAt.Time,
			"counter":    gorm.Expr("counter + 1"),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update refresh token: %v", result.Error)
	}

	// ตรวจสอบว่ามีแถวที่ถูกอัปเดตหรือไม่
	if result.RowsAffected == 0 {
		return fmt.Errorf("refresh token not found")
	}

	return nil
}

func (s *userRepository) RevokeRefreshToken(tokenID string) error {
	result := s.db.Model(&entities.RefreshToken{}).
		Where("id = ?", tokenID).
		Update("revoke", true)

	if result.Error != nil {
		return fmt.Errorf("failed to revoke refresh token: %v", result.Error)
	}

	// ตรวจสอบว่ามีแถวที่ถูกอัปเดตหรือไม่
	if result.RowsAffected == 0 {
		return fmt.Errorf("refresh token not found")
	}

	return nil
}
