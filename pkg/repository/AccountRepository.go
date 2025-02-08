package repository

// repo layer ที่ใช้เข้าถึง db

import (
	"go_income_outflow/entities"

	"gorm.io/gorm"
)

type (
	AccountRepository interface {
		Create(account *entities.Account) error
		FindAccountWithRelations(account *entities.Account, relations []string) error
		CreateWithRelations(account *entities.Account, relations []string) error
	}

	accountRepository struct {
		db *gorm.DB
	}
)

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(account *entities.Account) error {
	return r.db.Create(account).Error
}

func (r *accountRepository) CreateWithRelations(account *entities.Account, relations []string) error {
	if err := r.db.Create(account).Error; err != nil {
		return err
	}

	query := r.db
	for _, relation := range relations {
		query = query.Preload(relation)
	}

	if err := query.First(account).Error; err != nil {
		return err
	}

	return nil
}

func (r *accountRepository) FindAccountWithRelations(account *entities.Account, relations []string) error {
	query := r.db
	for _, relation := range relations {
		query = query.Preload(relation)
	}

	if err := query.First(account).Error; err != nil {
		return err
	}

	return nil
}
