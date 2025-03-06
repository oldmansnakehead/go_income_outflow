package service

import (
	"fmt"
	"go_income_outflow/entities"
	"go_income_outflow/pkg/model"
	"go_income_outflow/pkg/repository"
	"time"

	"github.com/shopspring/decimal"
)

type (
	CreditCardDebtServiceUseCase interface {
		CreateDebt(request *model.CreditCardDebtRequest) (debts *[]entities.CreditCardDebt, err error)
		calculateDueDate(dueDay uint) time.Time
	}

	creditCardDebtService struct {
		repo repository.CreditCardDebtRepository
	}
)

func NewCreditCardDebtService(repo repository.CreditCardDebtRepository) CreditCardDebtServiceUseCase {
	return &creditCardDebtService{repo: repo}
}

func (s *creditCardDebtService) CreateDebt(request *model.CreditCardDebtRequest) (debts *[]entities.CreditCardDebt, err error) {
	tx := s.repo.BeginTransaction()

	creditCard, err := s.repo.GetCreditCardByID(request.CreditCardID)
	if err != nil {
		return nil, err
	}

	dueDate := s.calculateDueDate(creditCard.DueDate)
	debtItems := make([]entities.CreditCardDebt, 0)

	for i := uint(0); i < uint(request.InstallmentCount); i++ {
		if i > 0 {
			// เพิ่ม DueDate ทีละ 1 เดือน ยกเว้นงวดแรก
			dueDate = dueDate.AddDate(0, 1, 0)
		}
		creditCardDebt := entities.CreditCardDebt{
			Description:  fmt.Sprintf("%s (Installment %d/%d)", request.Description, i+1, request.InstallmentCount),
			Amount:       request.Amount.Div(decimal.NewFromInt(int64(request.InstallmentCount))),
			DueDate:      dueDate.Format("2006-01-02"),
			CreditCardID: creditCard.ID,
		}

		if err := s.repo.CreateCreditCardDebt(&creditCardDebt, tx); err != nil {
			s.repo.RollbackTransaction(tx)
			return nil, err
		}

		debtItems = append(debtItems, creditCardDebt)
	}

	if err := s.repo.CommitTransaction(tx); err != nil {
		return nil, err
	}

	return &debtItems, nil
}

func (s *creditCardDebtService) calculateDueDate(dueDay uint) time.Time {
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
