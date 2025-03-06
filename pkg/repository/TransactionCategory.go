package repository

import (
	"go_income_outflow/entities"
	"go_income_outflow/helpers"
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
		FindWithFilters(filters map[string]interface{}) ([]model.TransactionCategoryResponse, error)
		FindByName(name string) (*entities.TransactionCategory, error)
		GetTransactionCategoryByID(categoryId uint) (*entities.TransactionCategory, error)
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
	if err := r.db.Model(transactionCategory).Omit("CreatedAt").Save(transactionCategory).Error; err != nil {
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

func (r *transactionCategoryRepository) FindWithFilters(filters map[string]interface{}) ([]model.TransactionCategoryResponse, error) {
	var items []entities.TransactionCategory
	query := r.db.Model(&entities.TransactionCategory{})

	if value, ok := filters["name"]; ok {
		query = helpers.WhereConditions(query, "name", value)
	}

	if value, ok := filters["type"]; ok {
		query = helpers.WhereConditions(query, "type", value)
	}

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
	item := entities.TransactionCategory{}
	if err := r.db.First(&item, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *transactionCategoryRepository) GetTransactionCategoryByID(categoryId uint) (*entities.TransactionCategory, error) {
	var category entities.TransactionCategory
	if err := r.db.First(&category, categoryId).Error; err != nil {
		return nil, err
	}
	return &category, nil
}
