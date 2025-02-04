package main

import (
	"go_income_outflow/db"
	"go_income_outflow/db/migrations"
	"go_income_outflow/routes"
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
		// godotenv.Load("../../.env")
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	db.ConnectDB()
	migrations.Migrate()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true

	r := gin.Default()
	r.Use(cors.New(corsConfig))

	routes.InitialRoute(r)

	r.Run(":" + os.Getenv("PORT"))
}
