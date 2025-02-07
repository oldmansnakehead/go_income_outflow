package repository

// repo layer ที่ใช้เข้าถึง db

import (
	"go_income_outflow/entities"

	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(account *entities.Account) error
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(account *entities.Account) error {
	return r.db.Create(account).Error
}
