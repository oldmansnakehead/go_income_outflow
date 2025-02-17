package repository

import (
	"go_income_outflow/entities"
	"go_income_outflow/pkg/model"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type (
	TransactionCategoryRepository interface {
		Create(transactionCategory *entities.TransactionCategory, relations []string) error
		FirstWithRelations(transactionCategory *entities.TransactionCategory, relations []string) error
		Update(transactionCategory *entities.TransactionCategory, relations []string) error
		Delete(transactionCategory *entities.TransactionCategory) error
		GetBaseQuery() *gorm.DB
		FindWithFilters(filters map[string]interface{}) ([]model.TransactionCategoryResponse, error)
		FindByName(name string) (*entities.TransactionCategory, error)
	}

	transactionCategoryRepository struct {
		db *gorm.DB
	}
)

func NewTransactionCategoryRepository(db *gorm.DB) TransactionCategoryRepository {
	return &transactionCategoryRepository{db: db}
}

func (r *transactionCategoryRepository) Create(transactionCategory *entities.TransactionCategory, relations []string) error {
	if err := r.db.Create(transactionCategory).Error; err != nil {
		return err
	}

	if len(relations) > 0 {
		if err := r.FirstWithRelations(transactionCategory, relations); err != nil {
			return err
		}
	}

	return nil
}

func (r *transactionCategoryRepository) FirstWithRelations(transactionCategory *entities.TransactionCategory, relations []string) error {
	query := r.db
	for _, relation := range relations {
		query = query.Preload(relation)
	}

	if err := query.First(transactionCategory).Error; err != nil {
		return err
	}

	return nil
}

func (r *transactionCategoryRepository) Update(transactionCategory *entities.TransactionCategory, relations []string) error {
	if err := r.db.Save(transactionCategory).Error; err != nil {
		return err
	}

	if len(relations) > 0 {
		if err := r.FirstWithRelations(transactionCategory, relations); err != nil {
			return err
		}
	}

	return nil
}

func (r *transactionCategoryRepository) Delete(transactionCategory *entities.TransactionCategory) error {
	if err := r.db.Delete(transactionCategory).Error; err != nil {
		return err
	}
	return nil
}

func (r *transactionCategoryRepository) GetBaseQuery() *gorm.DB {
	return r.db.Model(&entities.TransactionCategory{})
}

func (r *transactionCategoryRepository) FindWithFilters(filters map[string]interface{}) ([]model.TransactionCategoryResponse, error) {
	var items []entities.TransactionCategory
	query := r.db.Model(&entities.TransactionCategory{})

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	var response []model.TransactionCategoryResponse
	if err := copier.Copy(&response, &items); err != nil {
		return nil, nil
	}

	return response, nil
}

func (r *transactionCategoryRepository) FindByName(name string) (*entities.TransactionCategory, error) {
	var item entities.TransactionCategory
	result := r.db.First(&item, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}
