package service

// service หรือ usecase จะมีการนำข้อมูลขาเข้าและข้อมูลขาออก
// ขาออกเช่น AccountRepository
// handler (controller) -> usecase (service) -> repo

import (
	"go_income_outflow/entities"
	"go_income_outflow/pkg/repository"
)

// interface แทนการนำเข้าของข้อมูล
type AccountServiceUseCase interface {
	CreateAccount(account *entities.Account) error
}

type AccountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) CreateAccount(account *entities.Account) error {
	return s.repo.Create(account)
}
