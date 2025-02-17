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

	hasTable := db.Migrator().HasTable(&entities.Account{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entities.Account{}); err != nil {
			return err
		}
	}

	repo := repository.NewAccountRepository(db)
	service := service.NewAccountService(repo)

	for _, data := range items {
		existing, err := repo.FindByName(data.Name)
		if err == nil && existing != nil {
			fmt.Printf("Account with name %s already exists, skipping...\n", data.Name)
			continue
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			// ถ้ามีข้อผิดพลาดอื่น ๆ เกิดขึ้น
			return err
		}

		item := entities.Account{
			Name:   data.Name,
			UserID: data.UserID,
		}

		if err := service.CreateAccount(&item, []string{}); err != nil {
			return fmt.Errorf("failed to create account %s: %v", data.Name, err)
		}
		fmt.Printf("account %s created successfully\n", data.Name)
	}

	return nil
}
