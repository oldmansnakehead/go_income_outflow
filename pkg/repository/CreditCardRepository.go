package repository

import (
	"go_income_outflow/entities"
	"go_income_outflow/helpers"
	"go_income_outflow/pkg/model"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type (
	CreditCardRepository interface {
		Create(creditCard *entities.CreditCard, relations []string) error
		FirstWithRelations(creditCard *entities.CreditCard, relations []string) error
		Update(creditCard *entities.CreditCard, relations []string) error
		Delete(creditCard *entities.CreditCard) error
		FindWithFilters(filters map[string]interface{}) ([]model.CreditCardResponse, error)
	}

	creditCardRepository struct {
		db *gorm.DB
	}
)

func NewCreditCardRepository(db *gorm.DB) CreditCardRepository {
	return &creditCardRepository{db: db}
}

func (r *creditCardRepository) Create(creditCard *entities.CreditCard, relations []string) error {
	if err := r.db.Create(creditCard).Error; err != nil {
		return err
	}

	if len(relations) > 0 {
		if err := r.FirstWithRelations(creditCard, relations); err != nil {
			return err
		}
	}

	return nil
}

func (r *creditCardRepository) FirstWithRelations(creditCard *entities.CreditCard, relations []string) error {
	query := r.db
	for _, relation := range relations {
		query = query.Preload(relation)
	}

	if err := query.First(creditCard).Error; err != nil {
		return err
	}

	return nil
}

func (r *creditCardRepository) Update(creditCard *entities.CreditCard, relations []string) error {
	if err := r.db.Model(creditCard).Omit("CreatedAt").Save(creditCard).Error; err != nil {
		return err
	}

	if len(relations) > 0 {
		if err := r.FirstWithRelations(creditCard, relations); err != nil {
			return err
		}
	}

	return nil
}

func (r *creditCardRepository) Delete(creditCard *entities.CreditCard) error {
	if err := r.db.Delete(creditCard).Error; err != nil {
		return err
	}
	return nil
}

func (r *creditCardRepository) FindWithFilters(filters map[string]interface{}) ([]model.CreditCardResponse, error) {
	var creditCard []entities.CreditCard
	query := r.db.Model(&entities.CreditCard{})

	if relations, ok := filters["with"]; ok {
		query = helpers.WithRelations(query, relations)
	}

	if value, ok := filters["user_id"]; ok {
		query = helpers.WhereConditions(query, "user_id", value)
	}

	// ดึงข้อมูล
	if err := query.Find(&creditCard).Error; err != nil {
		return nil, err
	}
	var response []model.CreditCardResponse
	if err := copier.Copy(&response, &creditCard); err != nil {
		return nil, nil
	}

	return response, nil
}
