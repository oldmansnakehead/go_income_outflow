package repository

import (
	"go_income_outflow/entities"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *entities.User) error
	FindUserByEmail(email string) (*entities.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
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
