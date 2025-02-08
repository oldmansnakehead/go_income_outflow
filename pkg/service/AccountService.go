package service

// service หรือ usecase จะมีการนำข้อมูลขาเข้าและข้อมูลขาออก
// ขาออกเช่น AccountRepository
// handler (controller) -> usecase (service) -> repo

import (
	"go_income_outflow/entities"
	"go_income_outflow/pkg/repository"
)

type (
	// interface แทนการนำเข้าของข้อมูล
	AccountServiceUseCase interface {
		CreateAccount(account *entities.Account) error
		CreateAccountWithRelations(account *entities.Account, relations []string) error
		GetAccountWithRelations(account *entities.Account, relations []string) error
	}

	AccountService struct {
		repo repository.AccountRepository
	}
)

func NewAccountService(repo repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) CreateAccount(account *entities.Account) error {
	return s.repo.Create(account)
}

func (s *AccountService) CreateAccountWithRelations(account *entities.Account, relations []string) error {
	return s.repo.CreateWithRelations(account, relations)
}

func (s *AccountService) GetAccountWithRelations(account *entities.Account, relations []string) error {
	return s.repo.FindAccountWithRelations(account, relations)
}
