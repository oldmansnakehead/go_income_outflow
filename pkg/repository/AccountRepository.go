package repository

// repo layer ที่ใช้เข้าถึง db

import (
	"go_income_outflow/entities"
	"go_income_outflow/helpers"
	"go_income_outflow/pkg/model"

	"github.com/jinzhu/copier"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type (
	AccountRepository interface {
		Create(account *entities.Account, relations []string) error
		FirstWithRelations(account *entities.Account, relations []string) error
		Update(account *entities.Account, relations []string) error
		Delete(account *entities.Account) error
		GetBaseQuery() *gorm.DB
		FindWithFilters(filters map[string]interface{}) ([]model.AccountResponse, error)
		FindByName(name string) (*entities.Account, error)
		GetTotalAmount(userID uint) (decimal.Decimal, error)
	}

	accountRepository struct {
		db *gorm.DB
	}
)

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(account *entities.Account, relations []string) error {
	if err := r.db.Create(account).Error; err != nil {
		return err
	}

	if len(relations) > 0 {
		if err := r.FirstWithRelations(account, relations); err != nil {
			return err
		}
	}

	return nil
}

func (r *accountRepository) FirstWithRelations(account *entities.Account, relations []string) error {
	query := r.db
	for _, relation := range relations {
		query = query.Preload(relation)
	}

	if err := query.First(account).Error; err != nil {
		return err
	}

	return nil
}

func (r *accountRepository) Update(account *entities.Account, relations []string) error {
	// .Select("Name", "Amount", "ExcludeFromTotal", "Currency", "UserID") = อนุญาติให้อัพเดท fields ไหนบ้าง
	// Omit ไม่อัพเดท field ที่เลือก
	if err := r.db.Model(account).Omit("CreatedAt").Save(account).Error; err != nil {
		return err
	}

	if len(relations) > 0 {
		if err := r.FirstWithRelations(account, relations); err != nil {
			return err
		}
	}

	return nil
}

func (r *accountRepository) Delete(account *entities.Account) error {
	if err := r.db.Delete(account).Error; err != nil {
		return err
	}
	return nil
}

func (r *accountRepository) GetBaseQuery() *gorm.DB {
	return r.db.Model(&entities.Account{})
}

func (r *accountRepository) FindWithFilters(filters map[string]interface{}) ([]model.AccountResponse, error) {
	var accounts []entities.Account
	query := r.db.Model(&entities.Account{})

	if relations, ok := filters["with"]; ok {
		query = helpers.WithRelations(query, relations)

	}

	if value, ok := filters["user_id"]; ok {
		query = helpers.WhereConditions(query, "user_id", value)
	}

	if value, ok := filters["id"]; ok {
		query = helpers.WhereConditions(query, "id", value)
	}

	// ดึงข้อมูล
	if err := query.Find(&accounts).Error; err != nil {
		return nil, err
	}

	var response []model.AccountResponse
	if err := copier.Copy(&response, &accounts); err != nil {
		return nil, nil
	}

	return response, nil
}

func (r *accountRepository) FindByName(name string) (*entities.Account, error) {
	var item entities.Account
	result := r.db.First(&item, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *accountRepository) GetTotalAmount(userID uint) (decimal.Decimal, error) {
	var totalAmount decimal.Decimal

	// กรองตาม user_id
	err := r.db.Model(&entities.Account{}).
		Where("user_id = ?", userID). // กรองเฉพาะ user_id ที่ตรงกับที่ส่งมา
		Select("SUM(amount)").
		Scan(&totalAmount).Error

	if err != nil {
		return decimal.Zero, err
	}

	return totalAmount, nil
}
