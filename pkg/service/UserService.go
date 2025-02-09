package service

import (
	"fmt"
	"go_income_outflow/entities"
	"go_income_outflow/pkg/model"
	"go_income_outflow/pkg/repository"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type (
	UserServiceUseCase interface {
		CreateUser(body *model.UserRequest) error
		Login(email, password string) (string, error)
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

func (s *userService) Login(email, password string) (string, error) {
	// ค้นหาผู้ใช้จาก email
	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return "", fmt.Errorf("invalid email")
	}

	// เปรียบเทียบรหัสผ่าน
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid password")
	}

	// สร้าง token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// ลงนามและสร้าง token
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", fmt.Errorf("failed to create token")
	}

	return tokenString, nil
}
