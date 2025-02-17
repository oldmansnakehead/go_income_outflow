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

func TransactionCategorySeeder(db *gorm.DB) error {
	jsonFile, err := os.Open("../db/migrations/json/transaction_category.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, _ := io.ReadAll(jsonFile)

	var items []entities.TransactionCategory
	if err := json.Unmarshal(jsonData, &items); err != nil {
		return err
	}

	hasTable := db.Migrator().HasTable(&entities.TransactionCategory{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entities.TransactionCategory{}); err != nil {
			return err
		}
	}

	repo := repository.NewTransactionCategoryRepository(db)
	service := service.NewTransactionCategoryService(repo)

	for _, data := range items {
		existing, err := repo.FindByName(data.Name)
		if err == nil && existing != nil {
			fmt.Printf("category with name %s already exists, skipping...\n", data.Name)
			continue
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		item := entities.TransactionCategory{
			Name: data.Name,
			Type: data.Type,
		}

		if err := service.CreateTransactionCategory(&item, []string{}); err != nil {
			return fmt.Errorf("failed to create category %s: %v", data.Name, err)
		}
		fmt.Printf("category %s created successfully\n", data.Name)
	}

	return nil
}
