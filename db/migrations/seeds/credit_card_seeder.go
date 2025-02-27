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

func CreditCardSeeder(db *gorm.DB) error {
	jsonFile, err := os.Open("../db/migrations/json/credit_card.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	fmt.Println("test1")
	jsonData, _ := io.ReadAll(jsonFile)
	fmt.Println("test2")
	var items []entities.CreditCard
	if err := json.Unmarshal(jsonData, &items); err != nil {
		return err
	}
	fmt.Println("test3")
	hasTable := db.Migrator().HasTable(&entities.CreditCard{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entities.CreditCard{}); err != nil {
			return err
		}
	}

	repo := repository.NewCreditCardRepository(db)
	service := service.NewCreditCardService(repo)

	for _, data := range items {
		existing, err := repo.FindByName(data.Name)
		if err == nil && existing != nil {
			fmt.Printf("CreditCard with name %s already exists, skipping...\n", data.Name)
			continue
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		item := entities.CreditCard{
			Name:        data.Name,
			CreditLimit: data.CreditLimit,
			Balance:     data.Balance,
			DueDate:     data.DueDate,
			UserID:      data.UserID,
		}

		fmt.Println(item, "<-------------")

		if err := service.CreateCreditCard(&item, []string{}); err != nil {
			return fmt.Errorf("failed to create credit card %s: %v", data.Name, err)
		}
		fmt.Printf("credit card %s created successfully\n", data.Name)
	}

	return nil
}
