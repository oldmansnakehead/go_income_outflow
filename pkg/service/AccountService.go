package service

// service หรือ usecase จะมีการนำข้อมูลขาเข้าและข้อมูลขาออก
// ขาออกเช่น AccountRepository
// handler (controller) -> usecase (service) -> repo

import (
	"go_income_outflow/entities"
	"go_income_outflow/pkg/model"
	"go_income_outflow/pkg/repository"

	"github.com/shopspring/decimal"
)

type (
	// interface แทนการนำเข้าของข้อมูล
	AccountServiceUseCase interface {
		CreateAccount(account *entities.Account, relations []string) error
		FirstWithRelations(account *entities.Account, relations []string) error
		UpdateAccount(account *entities.Account, relations []string) error
		DeleteAccount(account *entities.Account) error
		GetWithFilters(filters map[string]interface{}) ([]model.AccountResponse, error)
		GetTotalAmount(userID uint) (decimal.Decimal, error)
	}

	accountService struct {
		repo repository.AccountRepository
	}
)

func NewAccountService(repo repository.AccountRepository) AccountServiceUseCase {
	return &accountService{repo: repo}
}

func (s *accountService) CreateAccount(account *entities.Account, relations []string) error {
	return s.repo.Create(account, relations)
}

func (s *accountService) FirstWithRelations(account *entities.Account, relations []string) error {
	return s.repo.FirstWithRelations(account, relations)
}

func (s *accountService) UpdateAccount(account *entities.Account, relations []string) error {
	return s.repo.Update(account, relations)
}

func (s *accountService) DeleteAccount(account *entities.Account) error {
	return s.repo.Delete(account)
}

func (s *accountService) GetWithFilters(filters map[string]interface{}) ([]model.AccountResponse, error) {
	return s.repo.FindWithFilters(filters)
}

func (s *accountService) GetTotalAmount(userID uint) (decimal.Decimal, error) {
	return s.repo.GetTotalAmount(userID)
}
