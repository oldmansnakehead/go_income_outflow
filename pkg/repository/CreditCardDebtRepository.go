package repository

import (
	"go_income_outflow/entities"

	"gorm.io/gorm"
)

type (
	CreditCardDebtRepository interface {
		BeginTransaction() *gorm.DB
		CommitTransaction(tx *gorm.DB) error
		RollbackTransaction(tx *gorm.DB) error
		CreateCreditCardDebt(creditCardDebt *entities.CreditCardDebt, tx *gorm.DB) error
		GetCreditCardByID(creditCardID uint) (*entities.CreditCard, error)
		GetCreditCardDeptByPaymentId(paymentID uint) (*entities.CreditCardDebt, error)
		UpdateCreditCardDebt(creditCardDebt *entities.CreditCardDebt, tx *gorm.DB) error
	}

	creditCardDebtRepository struct {
		db             *gorm.DB
		creditCardRepo CreditCardRepository
	}
)

func NewCreditCardDebtRepository(db *gorm.DB, creditCardRepo CreditCardRepository) CreditCardDebtRepository {
	return &creditCardDebtRepository{db: db, creditCardRepo: creditCardRepo}
}

func (r *creditCardDebtRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *creditCardDebtRepository) CommitTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (r *creditCardDebtRepository) RollbackTransaction(tx *gorm.DB) error {
	return tx.Rollback().Error
}

func (r *creditCardDebtRepository) CreateCreditCardDebt(creditCardDebt *entities.CreditCardDebt, tx *gorm.DB) error {
	conn := r.db
	if tx != nil {
		conn = tx
	}

	return conn.Create(creditCardDebt).Error
}

func (r *creditCardDebtRepository) GetCreditCardByID(creditCardID uint) (*entities.CreditCard, error) {
	return r.creditCardRepo.GetCreditCardByID(creditCardID)
}

func (r *creditCardDebtRepository) GetCreditCardDeptByPaymentId(paymentID uint) (*entities.CreditCardDebt, error) {
	var creditCardDebt entities.CreditCardDebt
	if err := r.db.Where("payment_transaction_id = ?", paymentID).First(&creditCardDebt).Error; err != nil {
		return nil, err
	}
	return &creditCardDebt, nil
}

func (r *creditCardDebtRepository) UpdateCreditCardDebt(creditCardDebt *entities.CreditCardDebt, tx *gorm.DB) error {
	conn := r.db
	if tx != nil {
		conn = tx
	}

	return conn.Save(creditCardDebt).Error
}
