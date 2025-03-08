package server

import (
	"go_income_outflow/db"
	"go_income_outflow/db/migrations"
	"go_income_outflow/routes"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Start() {
	if os.Getenv("APP_ENV") != "production" {
		// godotenv.Load("../../.env")
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	loc, err := time.LoadLocation(os.Getenv("TZ"))
	if err != nil {
		log.Fatalf("Error loading time zone: %v", err)
	}
	time.Local = loc // ตั้งค่าโซนเวลาให้กับแอปพลิเคชัน

	db := db.ConnectDB()
	migrations.Migrate()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true

	r := gin.Default()
	r.Use(cors.New(corsConfig))

	routes.InitialRoute(r, db)

	r.Run(":" + os.Getenv("PORT"))
}
