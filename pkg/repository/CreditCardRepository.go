package repository

import (
	"go_income_outflow/entities"

	"gorm.io/gorm"
)

type (
	CreditCardRepository interface {
		Create(creditCard *entities.CreditCard, relations []string) error
		FirstWithRelations(creditCard *entities.CreditCard, relations []string) error
		Update(creditCard *entities.CreditCard, relations []string) error
		Delete(creditCard *entities.CreditCard) error
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
	if err := r.db.Save(creditCard).Error; err != nil {
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
