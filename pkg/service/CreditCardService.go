package service

import (
	"go_income_outflow/entities"
	"go_income_outflow/pkg/repository"
)

type (
	CreditCardServiceUseCase interface {
		CreateCreditCard(creditCard *entities.CreditCard, relations []string) error
		FirstWithRelations(creditCard *entities.CreditCard, relations []string) error
		UpdateCreditCard(creditCard *entities.CreditCard, relations []string) error
		DeleteCreditCard(creditCard *entities.CreditCard) error
	}

	creditCardService struct {
		repo repository.CreditCardRepository
	}
)

func NewCreditCardService(repo repository.CreditCardRepository) CreditCardServiceUseCase {
	return &creditCardService{repo: repo}
}

func (s *creditCardService) CreateCreditCard(creditCard *entities.CreditCard, relations []string) error {
	return s.repo.Create(creditCard, relations)
}

func (s *creditCardService) FirstWithRelations(creditCard *entities.CreditCard, relations []string) error {
	return s.repo.FirstWithRelations(creditCard, relations)
}

func (s *creditCardService) UpdateCreditCard(creditCard *entities.CreditCard, relations []string) error {
	return s.repo.Update(creditCard, relations)
}

func (s *creditCardService) DeleteCreditCard(creditCard *entities.CreditCard) error {
	return s.repo.Delete(creditCard)
}
