package service

import (
	"go_income_outflow/entities"
	"go_income_outflow/pkg/model"
	"go_income_outflow/pkg/repository"
)

type (
	TransactionCategoryServiceUseCase interface {
		CreateTransactionCategory(transactionCategory *entities.TransactionCategory, relations []string) error
		FirstWithRelations(transactionCategory *entities.TransactionCategory, relations []string) error
		UpdateTransactionCategory(transactionCategory *entities.TransactionCategory, relations []string) error
		DeleteTransactionCategory(transactionCategory *entities.TransactionCategory) error
		GetWithFilters(filters map[string]interface{}) ([]model.TransactionCategoryResponse, error)
	}

	transactionCategoryService struct {
		repo repository.TransactionCategoryRepository
	}
)

func NewTransactionCategoryService(repo repository.TransactionCategoryRepository) TransactionCategoryServiceUseCase {
	return &transactionCategoryService{repo: repo}
}

func (s *transactionCategoryService) CreateTransactionCategory(transactionCategory *entities.TransactionCategory, relations []string) error {
	return s.repo.Create(transactionCategory, relations)
}

func (s *transactionCategoryService) FirstWithRelations(transactionCategory *entities.TransactionCategory, relations []string) error {
	return s.repo.FirstWithRelations(transactionCategory, relations)
}

func (s *transactionCategoryService) UpdateTransactionCategory(transactionCategory *entities.TransactionCategory, relations []string) error {
	return s.repo.Update(transactionCategory, relations)
}

func (s *transactionCategoryService) DeleteTransactionCategory(transactionCategory *entities.TransactionCategory) error {
	return s.repo.Delete(transactionCategory)
}

func (s *transactionCategoryService) GetWithFilters(filters map[string]interface{}) ([]model.TransactionCategoryResponse, error) {
	return s.repo.FindWithFilters(filters)
}
