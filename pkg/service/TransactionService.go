package service

import (
	"errors"
	"go_income_outflow/entities"
	"go_income_outflow/pkg/model"
	"go_income_outflow/pkg/repository"
	"time"

	"gorm.io/gorm"
)

type (
	TransactionServiceUseCase interface {
		CreateTransaction(transaction *entities.Transaction, request model.TransactionRequest) error
		FirstWithRelations(transaction *entities.Transaction, relations []string) error
		DeleteTransaction(transaction *entities.Transaction) error
		GetWithFilters(filters map[string]interface{}) ([]model.TransactionResponse, error)
		calculateDueDate(dueDay uint) time.Time
		handleAccountTransaction(transaction *entities.Transaction, tx *gorm.DB) error
		handleCreditCardTransaction(transaction *entities.Transaction, request model.TransactionRequest, tx *gorm.DB) error
		handleCreditCardDept(transaction *entities.Transaction, creditCardDeptID uint, tx *gorm.DB) error
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
		if err := s.handleAccountTransaction(transaction, tx); err != nil {
			s.repo.RollbackTransaction(tx)
			return err
		}

		if request.CreditCardDebtId != 0 {
			if err := s.handleCreditCardDept(transaction, request.CreditCardDebtId, tx); err != nil {
				s.repo.RollbackTransaction(tx)
				return err
			}
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
	tx := s.repo.BeginTransaction()

	transaction, err := s.repo.GetTransactionByID(transaction.ID)
	if err != nil {
		return err
	}

	if transaction.TransactionableType == "accounts" {
		creditCardDept, err := s.repo.GetCreditCardDeptByPaymentId(transaction.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if creditCardDept != nil {
			creditCardDept.PaymentTransactionID = nil
			if err := s.repo.UpdateCreditCardDebt(creditCardDept, tx); err != nil {
				tx.Rollback()
				return err
			}
		}

		account, err := s.repo.GetAccountByID(transaction.TransactionableID)
		if err != nil {
			return err
		}

		account.Balance = account.Balance.Add(transaction.Amount.Neg())
		if err := s.repo.UpdateAccountBalance(account, tx); err != nil {
			tx.Rollback()
			return err
		}
	} else if transaction.TransactionableType == "credit_cards" {
		creditCard, err := s.repo.GetCreditCardByID(transaction.TransactionableID)
		if err != nil {
			return err
		}
		creditCard.Balance = creditCard.Balance.Add(transaction.Amount.Neg())
		if err := s.repo.UpdateCreditCardBalance(creditCard, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := s.repo.Delete(transaction, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.CommitTransaction(tx); err != nil {
		return err
	}

	return nil
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

	creditCardDebt := entities.CreditCardDebt{
		Description:   transaction.Description,
		Amount:        transaction.Amount,
		DueDate:       dueDate.Format("2006-01-02"),
		CreditCardID:  creditCard.ID,
		TransactionID: &transaction.ID,
	}

	if err := s.repo.CreateCreditCardDebt(&creditCardDebt, tx); err != nil {
		return err
	}

	return nil
}

func (s *transactionService) handleCreditCardDept(transaction *entities.Transaction, creditCardDeptID uint, tx *gorm.DB) error {
	creditCardDept, err := s.repo.GetCreditCardDeptByID(creditCardDeptID)
	if err != nil {
		return err
	}

	if creditCardDept.PaymentTransactionID != nil {
		return errors.New("payment has already been made for this credit card debt")
	}

	if transaction.Amount.Abs().LessThan(creditCardDept.Amount.Abs()) {
		return errors.New("payment amount is insufficient to cover the credit card debt")
	}

	creditCardDept.Amount = transaction.Amount
	creditCardDept.PaymentTransactionID = &transaction.ID

	if err := s.repo.UpdateCreditCardDebt(creditCardDept, tx); err != nil {
		return err
	}

	creditCard, err := s.repo.GetCreditCardByID(creditCardDept.CreditCardID)
	if err != nil {
		return err
	}

	creditCard.Balance = creditCard.Balance.Add(transaction.Amount.Abs())

	return s.repo.UpdateCreditCardBalance(creditCard, tx)
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
