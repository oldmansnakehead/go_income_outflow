package repository

import (
	"fmt"
	"go_income_outflow/entities"
	"go_income_outflow/helpers"
	"go_income_outflow/pkg/model"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type (
	TransactionRepository interface {
		BeginTransaction() *gorm.DB
		CommitTransaction(tx *gorm.DB) error
		RollbackTransaction(tx *gorm.DB) error

		CreateTransaction(transaction *entities.Transaction, tx *gorm.DB) error
		GetTransactionByID(transactionID uint) (*entities.Transaction, error)
		FirstWithRelations(transaction *entities.Transaction, relations []string) error
		Delete(transaction *entities.Transaction, tx *gorm.DB) error
		FindWithFilters(filters map[string]interface{}) ([]model.TransactionResponse, error)

		GetAccountByID(accountID uint) (*entities.Account, error)
		GetCreditCardDeptByID(CreditCardDebtID uint) (*entities.CreditCardDebt, error)
		GetCreditCardDeptByPaymentId(PaymentID uint) (*entities.CreditCardDebt, error)
		GetTransactionCategoryByID(categoryId uint) (*entities.TransactionCategory, error)
		GetCreditCardByID(creditCardID uint) (*entities.CreditCard, error)

		UpdateAccountBalance(account *entities.Account, tx *gorm.DB) error
		UpdateCreditCardBalance(creditCard *entities.CreditCard, tx *gorm.DB) error
		UpdateCreditCardDebt(creditCardDebt *entities.CreditCardDebt, tx *gorm.DB) error
		CreateCreditCardDebt(creditCardDebt *entities.CreditCardDebt, tx *gorm.DB) error
	}

	transactionRepository struct {
		db                      *gorm.DB
		accountRepo             AccountRepository
		creditCardRepo          CreditCardRepository
		creditCardDebtRepo      CreditCardDebtRepository
		transactionCategoryRepo TransactionCategoryRepository
	}
)

func NewTransactionRepository(
	db *gorm.DB,
	accountRepo AccountRepository,
	creditCardRepo CreditCardRepository,
	creditCardDebtRepo CreditCardDebtRepository,
	transactionCategoryRepo TransactionCategoryRepository,
) TransactionRepository {
	return &transactionRepository{
		db:                      db,
		accountRepo:             accountRepo,
		creditCardRepo:          creditCardRepo,
		creditCardDebtRepo:      creditCardDebtRepo,
		transactionCategoryRepo: transactionCategoryRepo,
	}
}

func (r *transactionRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *transactionRepository) CommitTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (r *transactionRepository) RollbackTransaction(tx *gorm.DB) error {
	return tx.Rollback().Error
}

func (r *transactionRepository) CreateTransaction(transaction *entities.Transaction, tx *gorm.DB) error {
	conn := r.db
	if tx != nil {
		conn = tx
	}
	return conn.Create(transaction).Error
}

func (r *transactionRepository) GetAccountByID(accountID uint) (*entities.Account, error) {
	return r.accountRepo.GetAccountByID(accountID)
}

func (r *transactionRepository) GetCreditCardDeptByID(CreditCardDebtID uint) (*entities.CreditCardDebt, error) {
	var creditCardDebt entities.CreditCardDebt
	if err := r.db.First(&creditCardDebt, CreditCardDebtID).Error; err != nil {
		return nil, err
	}
	return &creditCardDebt, nil
}

func (r *transactionRepository) GetCreditCardDeptByPaymentId(paymentID uint) (*entities.CreditCardDebt, error) {
	return r.creditCardDebtRepo.GetCreditCardDeptByPaymentId(paymentID)
}

func (r *transactionRepository) GetTransactionCategoryByID(categoryId uint) (*entities.TransactionCategory, error) {
	return r.transactionCategoryRepo.GetTransactionCategoryByID(categoryId)
}

func (r *transactionRepository) GetCreditCardByID(creditCardID uint) (*entities.CreditCard, error) {
	return r.creditCardRepo.GetCreditCardByID(creditCardID)
}

func (r *transactionRepository) GetTransactionByID(transactionID uint) (*entities.Transaction, error) {
	var transaction entities.Transaction
	if err := r.db.First(&transaction, transactionID).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) UpdateAccountBalance(account *entities.Account, tx *gorm.DB) error {
	return r.accountRepo.UpdateAccountBalance(account, tx)
}

func (r *transactionRepository) UpdateCreditCardBalance(creditCard *entities.CreditCard, tx *gorm.DB) error {
	return r.creditCardRepo.UpdateCreditCardBalance(creditCard, tx)
}

func (r *transactionRepository) CreateCreditCardDebt(creditCardDebt *entities.CreditCardDebt, tx *gorm.DB) error {
	return r.creditCardDebtRepo.CreateCreditCardDebt(creditCardDebt, tx)
}

func (r *transactionRepository) UpdateCreditCardDebt(creditCardDebt *entities.CreditCardDebt, tx *gorm.DB) error {
	return r.creditCardDebtRepo.UpdateCreditCardDebt(creditCardDebt, tx)
}

func (r *transactionRepository) FirstWithRelations(transaction *entities.Transaction, relations []string) error {
	query := r.db

	for _, relation := range relations {
		query = query.Preload(relation)
	}

	if err := query.First(transaction).Error; err != nil {
		return err
	}

	// preload morph (auto load)
	switch transaction.TransactionableType {
	case "credit_cards":
		var creditCard entities.CreditCard
		err := r.db.First(&creditCard, transaction.TransactionableID).Error
		if err != nil {
			return err
		}
		transaction.Transactionable = creditCard
	case "accounts":
		var account entities.Account
		err := r.db.First(&account, transaction.TransactionableID).Error
		if err != nil {
			return err
		}
		transaction.Transactionable = account
	default:
		return fmt.Errorf("unknown morphed type: %s", transaction.TransactionableType)
	}

	return nil
}

func (r *transactionRepository) Delete(transaction *entities.Transaction, tx *gorm.DB) error {
	conn := r.db
	if tx != nil {
		conn = tx
	}

	if err := conn.Delete(transaction).Error; err != nil {
		return err
	}
	return nil
}

func (r *transactionRepository) FindWithFilters(filters map[string]interface{}) ([]model.TransactionResponse, error) {
	var transactions []entities.Transaction
	query := r.db.Model(&entities.Transaction{})

	if relations, ok := filters["with"]; ok {
		query = helpers.WithRelations(query, relations)
	}

	if value, ok := filters["user_id"]; ok {
		query = helpers.WhereConditions(query, "user_id", value)
	}

	if value, ok := filters["category_id"]; ok {
		query = helpers.WhereConditions(query, "category_id", value)
	}

	if value, ok := filters["transactionable_id"]; ok {
		query = helpers.WhereConditions(query, "transactionable_id", value)
	}

	if value, ok := filters["transactionable_type"]; ok {
		query = helpers.WhereConditions(query, "transactionable_type", value)
	}

	if err := query.Find(&transactions).Error; err != nil {
		return nil, err
	}

	var response []model.TransactionResponse
	if err := copier.Copy(&response, &transactions); err != nil {
		return nil, nil
	}

	return response, nil
}
