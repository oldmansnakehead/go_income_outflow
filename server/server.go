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

	/* corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true */

	corsConfig := cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000", // สำหรับ development
		},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"auth.refresh-token",
			"auth.token",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // จำเป็นหากใช้ cookies หรือ authentication headers
		MaxAge:           12 * time.Hour,
	}

	/* corsConfig := cors.Config{
		AllowOrigins:     []string{"*"}, // ไม่สามารถใช้ร่วมกับ AllowCredentials = true
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "RefreshToken"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false, // ต้องเป็น false เมื่อใช้ AllowOrigins = ["*"]
		MaxAge:           12 * time.Hour,
	} */

	r := gin.Default()
	r.Use(cors.New(corsConfig))

	routes.InitialRoute(r, db)

	r.Run(":" + os.Getenv("PORT"))
}
