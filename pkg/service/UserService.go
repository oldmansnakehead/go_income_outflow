package service

import (
	"fmt"
	"go_income_outflow/entities"
	"go_income_outflow/pkg/model"
	"go_income_outflow/pkg/repository"
	"go_income_outflow/token"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type (
	UserServiceUseCase interface {
		CreateUser(body *model.UserRequest) error
		Login(email, password string) (refreshToken string, accessToken string, refreshExpiresAt time.Time, accessExpiresAt time.Time, err error)
		Logout(userID uint) error
		GetRefreshTokenByID(tokenID string) (*entities.RefreshToken, error)
		UpdateRefreshToken(rc *token.RefreshClaims, refreshToken string) error
		RevokeRefreshToken(tokenID string) error
	}

	userService struct {
		repo repository.UserRepository
	}
)

func NewUserService(repo repository.UserRepository) UserServiceUseCase {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(body *model.UserRequest) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	user := entities.User{Email: body.Email, Password: string(hash), Name: body.Name}
	err = s.repo.CreateUser(&user)

	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

func (s *userService) Login(email, password string) (refreshToken string, accessToken string, refreshExpiresAt time.Time, accessExpiresAt time.Time, err error) {
	tx := s.repo.BeginTransaction()

	// ค้นหาผู้ใช้จาก email
	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, fmt.Errorf("invalid email")
	}

	// เปรียบเทียบรหัสผ่าน
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", time.Time{}, time.Time{}, fmt.Errorf("invalid password")
	}

	u := &model.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	rc := token.NewRefreshClaims(u)
	ac := token.NewAccessClaims(u)

	refreshTokenString, err := rc.JwtString()
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	accessTokenString, err := ac.JwtString()
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	if err := s.repo.ClearRefreshToken(user.ID, tx); err != nil {
		s.repo.RollbackTransaction(tx)
		return "", "", time.Time{}, time.Time{}, err
	}

	expiresAt := rc.ExpiresAt.Time
	if err := s.repo.SaveRefreshToken(user.ID, refreshTokenString, expiresAt, rc.TokenID, tx); err != nil {
		s.repo.RollbackTransaction(tx)
		return "", "", time.Time{}, time.Time{}, fmt.Errorf("failed to save refresh token: %v", err)
	}

	if err := s.repo.CommitTransaction(tx); err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	return refreshTokenString, accessTokenString, rc.ExpiresAt.Time, ac.ExpiresAt.Time, nil
}

func (s *userService) Logout(userID uint) error {
	if err := s.repo.ClearRefreshToken(userID, nil); err != nil {
		return fmt.Errorf("failed to clear refresh tokens: %v", err)
	}

	return nil
}

func (s *userService) GetRefreshTokenByID(tokenID string) (*entities.RefreshToken, error) {
	return s.repo.GetRefreshTokenByID(tokenID)
}

func (s *userService) UpdateRefreshToken(rc *token.RefreshClaims, refreshToken string) error {
	return s.repo.UpdateRefreshToken(rc, refreshToken)
}

func (s *userService) RevokeRefreshToken(tokenID string) error {
	return s.repo.RevokeRefreshToken(tokenID)
}
