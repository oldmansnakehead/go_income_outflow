package main

import (
	"go_income_outflow/databases"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var Conn *gorm.DB

func main() {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load("../../.env"); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	databases.ConnectDB()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true

	r := gin.Default()
	r.Use(cors.New(corsConfig))
}
