package service

import (
	"errors"
	"fmt"
	"go_income_outflow/entities"
	"go_income_outflow/pkg/model"
	"go_income_outflow/pkg/repository"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type (
	TransactionServiceUseCase interface {
		CreateTransaction(transaction *entities.Transaction, request model.TransactionRequest) error
		FirstWithRelations(transaction *entities.Transaction, relations []string) error
		DeleteTransaction(transaction *entities.Transaction) error
		GetWithFilters(filters map[string]interface{}) ([]model.TransactionResponse, error)

		handleAccountTransaction(transaction *entities.Transaction, tx *gorm.DB) error
		handleCreditCardTransaction(transaction *entities.Transaction, request model.TransactionRequest, tx *gorm.DB) error
		calculateDueDate(dueDay uint) time.Time
	}

	transactionService struct {
		repo repository.TransactionRepository
	}
)

func NewTransactionService(repo repository.TransactionRepository) TransactionServiceUseCase {
	return &transactionService{repo: repo}
}

func (s *transactionService) CreateTransaction(transaction *entities.Transaction, request model.TransactionRequest) error {
	tx := s.repo.BeginTransaction()

	category, err := s.repo.GetTransactionCategoryByID(transaction.CategoryID)
	if err != nil {
		return err
	}

	// จัดการยอด (ส่งยอดมาเป็น -/+ ได้หมด)
	if !category.Type && transaction.Amount.Sign() > 0 || category.Type && transaction.Amount.Sign() < 0 {
		transaction.Amount = transaction.Amount.Neg() // .Neg() = * -1
	}

	if err := s.repo.CreateTransaction(transaction, tx); err != nil {
		s.repo.RollbackTransaction(tx)
		return err
	}

	if transaction.TransactionableType == "accounts" {
		err := s.handleAccountTransaction(transaction, tx)
		if err != nil {
			s.repo.RollbackTransaction(tx)
			return err
		}
	}

	if transaction.TransactionableType == "credit_cards" {
		err := s.handleCreditCardTransaction(transaction, request, tx)
		if err != nil {
			s.repo.RollbackTransaction(tx)
			return err
		}
	}

	if err := s.repo.CommitTransaction(tx); err != nil {
		return err
	}

	return nil
}

func (s *transactionService) FirstWithRelations(transaction *entities.Transaction, relations []string) error {
	return s.repo.FirstWithRelations(transaction, relations)
}

func (s *transactionService) DeleteTransaction(transaction *entities.Transaction) error {
	return s.repo.Delete(transaction)
}

func (s *transactionService) GetWithFilters(filters map[string]interface{}) ([]model.TransactionResponse, error) {
	return s.repo.FindWithFilters(filters)
}

func (s *transactionService) handleAccountTransaction(transaction *entities.Transaction, tx *gorm.DB) error {
	account, err := s.repo.GetAccountByID(transaction.TransactionableID)
	if err != nil {
		return err
	}

	if account.Balance.Add(transaction.Amount).Sign() < 0 {
		return errors.New("insufficient funds: transaction cannot be completed")
	}

	account.Balance = account.Balance.Add(transaction.Amount)
	return s.repo.UpdateAccountBalance(account, tx)
}

func (s *transactionService) handleCreditCardTransaction(transaction *entities.Transaction, request model.TransactionRequest, tx *gorm.DB) error {
	creditCard, err := s.repo.GetCreditCardByID(transaction.TransactionableID)
	if err != nil {
		return err
	}
	if creditCard.Balance.Add(transaction.Amount).Sign() < 0 {
		return errors.New("insufficient funds: transaction cannot be completed")
	}
	creditCard.Balance = creditCard.Balance.Add(transaction.Amount)
	if err := s.repo.UpdateCreditCardBalance(creditCard, tx); err != nil {
		return err
	}

	dueDate := s.calculateDueDate(creditCard.DueDate)

	if request.InstallmentCount > 0 {
		// วนลูปสร้าง CreditCardDebt ตามจำนวน Installments
		for i := uint(0); i < uint(request.InstallmentCount); i++ {
			if i > 0 {
				// เพิ่ม DueDate ทีละ 1 เดือน ยกเว้นงวดแรก
				dueDate = dueDate.AddDate(0, 1, 0)
			}
			creditCardDebt := entities.CreditCardDebt{
				Description:   fmt.Sprintf("%s (Installment %d/%d)", transaction.Description, i+1, request.InstallmentCount),
				Amount:        transaction.Amount.Div(decimal.NewFromInt(int64(request.InstallmentCount))),
				DueDate:       dueDate.Format("2006-01-02"),
				CreditCardID:  creditCard.ID,
				TransactionID: transaction.ID,
			}

			if err := s.repo.CreateCreditCardDebt(&creditCardDebt, tx); err != nil {
				return err
			}
		}
	} else {
		creditCardDebt := entities.CreditCardDebt{
			Description:   transaction.Description,
			Amount:        transaction.Amount,
			DueDate:       dueDate.Format("2006-01-02"),
			CreditCardID:  creditCard.ID,
			TransactionID: transaction.ID,
		}

		if err := s.repo.CreateCreditCardDebt(&creditCardDebt, tx); err != nil {
			return err
		}
	}

	return nil
}

func (s *transactionService) calculateDueDate(dueDay uint) time.Time {
	now := time.Now()
	year, month, day := now.Date()
	loc := now.Location()

	if uint(day) > dueDay {
		// ถ้าข้าม DueDate ไปแล้ว ใช้เดือนถัดไป
		month++
		if month > 12 {
			month = 1
			year++
		}
	}

	dueDate := time.Date(year, month, int(dueDay), 0, 0, 0, 0, loc)
	return dueDate // แปลงเป็น string
}
