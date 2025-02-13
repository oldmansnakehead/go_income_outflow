package seeds

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_income_outflow/entities"
	"go_income_outflow/pkg/repository"
	"go_income_outflow/pkg/service"
	"io"
	"os"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func AccountSeeder(db *gorm.DB) error {
	jsonFile, err := os.Open("../db/migrations/json/account.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, _ := io.ReadAll(jsonFile)

	var items []entities.Account
	if err := json.Unmarshal(jsonData, &items); err != nil {
		return err
	}

	hasTable := db.Migrator().HasTable(&entities.User{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entities.User{}); err != nil {
			return err
		}
	}

	repo := repository.NewAccountRepository(db)
	service := service.NewAccountService(repo)

	for _, data := range items {
		existing, err := repo.FindByName(data.Name)
		if err == nil && existing != nil {
			// ถ้ามีผู้ใช้อยู่แล้ว ไม่ต้องทำการสร้างใหม่
			fmt.Printf("User with email %s already exists, skipping...\n", data.Name)
			continue
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			// ถ้ามีข้อผิดพลาดอื่น ๆ เกิดขึ้น
			return err
		}

		var item entities.Account
		if err := copier.Copy(&item, &data); err != nil {
			return nil
		}

		if err := service.CreateAccount(&item, []string{}); err != nil {
			return fmt.Errorf("failed to create account %s: %v", data.Name, err)
		}
		fmt.Printf("User %s created successfully\n", data.Name)
	}

	return nil
}
