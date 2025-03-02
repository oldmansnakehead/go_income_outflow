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
		GetAccountByID(accountID uint) (*entities.Account, error)
		GetTransactionCategoryByID(categoryId uint) (*entities.TransactionCategory, error)
		GetCreditCardByID(creditCardID uint) (*entities.CreditCard, error)
		UpdateAccountBalance(account *entities.Account, tx *gorm.DB) error
		UpdateCreditCardBalance(creditCard *entities.CreditCard, tx *gorm.DB) error
		CreateCreditCardDebt(creditCardDebt *entities.CreditCardDebt, tx *gorm.DB) error
		Create(transaction *entities.Transaction, relations []string, request model.TransactionRequest) error
		FirstWithRelations(transaction *entities.Transaction, relations []string) error
		Delete(transaction *entities.Transaction) error
		FindWithFilters(filters map[string]interface{}) ([]model.TransactionResponse, error)
	}

	transactionRepository struct {
		db *gorm.DB
	}
)

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
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
	var account entities.Account
	if err := r.db.First(&account, accountID).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *transactionRepository) GetTransactionCategoryByID(categoryId uint) (*entities.TransactionCategory, error) {
	var account entities.TransactionCategory
	if err := r.db.First(&account, categoryId).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *transactionRepository) GetCreditCardByID(creditCardID uint) (*entities.CreditCard, error) {
	var creditCard entities.CreditCard
	if err := r.db.First(&creditCard, creditCardID).Error; err != nil {
		return nil, err
	}
	return &creditCard, nil
}

func (r *transactionRepository) UpdateAccountBalance(account *entities.Account, tx *gorm.DB) error {
	conn := r.db
	if tx != nil {
		conn = tx
	}

	return conn.Save(account).Error
}

func (r *transactionRepository) UpdateCreditCardBalance(creditCard *entities.CreditCard, tx *gorm.DB) error {
	conn := r.db
	if tx != nil {
		conn = tx
	}

	return conn.Save(creditCard).Error
}

func (r *transactionRepository) CreateCreditCardDebt(creditCardDebt *entities.CreditCardDebt, tx *gorm.DB) error {
	conn := r.db
	if tx != nil {
		conn = tx
	}

	return conn.Create(creditCardDebt).Error
}

func (r *transactionRepository) Create(transaction *entities.Transaction, relations []string, request model.TransactionRequest) error {
	tx := r.db.Begin()

	if err := r.db.Create(transaction).Error; err != nil {
		tx.Rollback()
		return err
	}

	if transaction.TransactionableType == "accounts" {
		account := entities.Account{}
		if err := r.db.First(&account, transaction.TransactionableID).Error; err != nil {
			tx.Rollback()
			return err
		}

		account.Balance = account.Balance.Add(transaction.Amount)
		if err := r.db.Save(&account).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if transaction.TransactionableType == "credit_cards" {
		var creditCard entities.CreditCard
		if err := r.db.First(&creditCard, transaction.TransactionableID).Error; err != nil {
			tx.Rollback()
			return err
		}

		creditCard.Balance = creditCard.Balance.Add(transaction.Amount)
		if err := r.db.Save(&creditCard).Error; err != nil {
			tx.Rollback()
			return err
		}

		if request.InstallmentCount > 0 {
			for i := uint(0); i < request.InstallmentCount; i++ {

			}
		} else {
			// dutDate := creditCard.DueDate

			creditCardDebt := entities.CreditCardDebt{
				Description: transaction.Description,
				Amount:      transaction.Amount,
				// DueDate:       dueDate,
				CreditCardID:  creditCard.ID,
				TransactionID: transaction.ID,
			}
			if err := r.db.Create(&creditCardDebt).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

	}

	if len(relations) > 0 {
		if err := r.FirstWithRelations(transaction, relations); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
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

func (r *transactionRepository) Delete(transaction *entities.Transaction) error {
	if err := r.db.Delete(transaction).Error; err != nil {
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
