package db

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Conn *gorm.DB
	once sync.Once
)

func ConnectDB() (c *gorm.DB) {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s search_path=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_SSL_MODE"),
			os.Getenv("DB_SCHEMA"),
		)

		db, err := gorm.Open(
			postgres.Open(dsn),
			&gorm.Config{Logger: logger.Default.LogMode(logger.Info)},
		)
		if err != nil {
			log.Fatal("Cannot connect to the database")
		}

		Conn = db
	})

	return Conn
}

// ใช้ใน cobra
func LazyConnect() *gorm.DB {
	// โหลดไฟล์ .env
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	return ConnectDB()
}
