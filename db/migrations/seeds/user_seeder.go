package seeds

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_income_outflow/entities"
	"go_income_outflow/pkg/model"
	"go_income_outflow/pkg/repository"
	"go_income_outflow/pkg/service"
	"io"
	"os"

	"gorm.io/gorm"
)

func UserSeeder(db *gorm.DB) error {
	jsonFile, err := os.Open("../db/migrations/json/user.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, _ := io.ReadAll(jsonFile)

	var items []entities.User
	if err := json.Unmarshal(jsonData, &items); err != nil {
		return err
	}

	// สร้างตารางใหม่ถ้าไม่มี
	hasTable := db.Migrator().HasTable(&entities.User{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entities.User{}); err != nil {
			return err
		}
	}

	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	repo := repository.NewUserRepository(db, refreshTokenRepo)
	userService := service.NewUserService(repo)

	for _, data := range items {
		existingUser, err := repo.FindUserByEmail(data.Email)
		if err == nil && existingUser != nil {
			// ถ้ามีผู้ใช้อยู่แล้ว ไม่ต้องทำการสร้างใหม่
			fmt.Printf("User with email %s already exists, skipping...\n", data.Email)
			continue
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			// ถ้ามีข้อผิดพลาดอื่น ๆ เกิดขึ้น
			return err
		}

		userRequest := model.UserRequest{
			Email:    data.Email,
			Password: data.Password,
			Name:     data.Name,
		}

		if err := userService.CreateUser(&userRequest); err != nil {
			return fmt.Errorf("failed to create user %s: %v", data.Email, err)
		}
		fmt.Printf("User %s created successfully\n", data.Email)
	}

	return nil
}
