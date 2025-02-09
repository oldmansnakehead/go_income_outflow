package repository

// repo layer ที่ใช้เข้าถึง db

import (
	"go_income_outflow/entities"
	"go_income_outflow/helpers"
	"go_income_outflow/pkg/model"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type (
	AccountRepository interface {
		Create(account *entities.Account, relations []string) error
		FirstWithRelations(account *entities.Account, relations []string) error
		Update(account *entities.Account, relations []string) error
		Delete(account *entities.Account) error
		GetBaseQuery() *gorm.DB
		FindWithFilters(filters map[string]interface{}) ([]model.AccountResponse, error)
	}

	accountRepository struct {
		db *gorm.DB
	}
)

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(account *entities.Account, relations []string) error {
	if err := r.db.Create(account).Error; err != nil {
		return err
	}

	if len(relations) > 0 {
		if err := r.FirstWithRelations(account, relations); err != nil {
			return err
		}
	}

	return nil
}

func (r *accountRepository) FirstWithRelations(account *entities.Account, relations []string) error {
	query := r.db
	for _, relation := range relations {
		query = query.Preload(relation)
	}

	if err := query.First(account).Error; err != nil {
		return err
	}

	return nil
}

func (r *accountRepository) Update(account *entities.Account, relations []string) error {
	if err := r.db.Save(account).Error; err != nil {
		return err
	}

	if len(relations) > 0 {
		if err := r.FirstWithRelations(account, relations); err != nil {
			return err
		}
	}

	return nil
}

func (r *accountRepository) Delete(account *entities.Account) error {
	if err := r.db.Delete(account).Error; err != nil {
		return err
	}
	return nil
}

func (r *accountRepository) GetBaseQuery() *gorm.DB {
	return r.db.Model(&entities.Account{})
}

func (r *accountRepository) FindWithFilters(filters map[string]interface{}) ([]model.AccountResponse, error) {
	var accounts []entities.Account
	query := r.db.Model(&entities.Account{})

	if relations, ok := filters["with"]; ok {
		query = helpers.WithRelations(query, relations)
	}

	if field, ok := filters["user_id"]; ok {
		query = helpers.SingleOrMultiple(query, "user_id", field)
	}

	// ดึงข้อมูล
	if err := query.Find(&accounts).Error; err != nil {
		return nil, err
	}
	var response []model.AccountResponse
	if err := copier.Copy(&response, &accounts); err != nil {
		return nil, nil
	}

	return response, nil
}
