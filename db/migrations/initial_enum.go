package migrations

import (
	"fmt"
	"go_income_outflow/db"
	"log"
	"strings"
)

// initialEnum ใช้สำหรับสร้างหลายๆ ENUM type ตามชื่อที่ส่งเข้ามา
func initialEnum(enumNames []string) {
	for _, enumName := range enumNames {
		// สร้าง query เพื่อเช็คว่า ENUM type มีอยู่หรือไม่
		var count int
		query := fmt.Sprintf(`SELECT COUNT(*) FROM pg_type WHERE typname = '%s'`, enumName)

		// ตรวจสอบว่า ENUM type มีอยู่ในฐานข้อมูลหรือไม่
		if err := db.Conn.Raw(query).Scan(&count).Error; err != nil {
			log.Fatal("Error checking for existing ENUM type:", err)
		}

		// ถ้า ENUM type ยังไม่ถูกสร้างขึ้น, ให้สร้างใหม่
		if count == 0 {
			// สร้างคำสั่ง SQL สำหรับสร้าง ENUM type
			enumValues := getEnumValues(enumName) // ฟังก์ชันที่จะดึงค่าของ ENUM ที่ต้องการ
			createEnumQuery := fmt.Sprintf("CREATE TYPE %s AS ENUM ('%s')", enumName, strings.Join(enumValues, "', '"))
			if err := db.Conn.Exec(createEnumQuery).Error; err != nil {
				log.Fatal("Failed to create ENUM type:", err)
			} else {
				fmt.Printf("Successfully created ENUM type '%s'\n", enumName)
			}
		} else {
			fmt.Printf("ENUM type '%s' already exists, skipping creation\n", enumName)
		}
	}
}

// getEnumValues คือตัวอย่างของฟังก์ชันที่ดึงค่าของ ENUM type
// ในกรณีนี้ คุณจะต้องจัดการกับค่า ENUM ของแต่ละ type โดยตรง
func getEnumValues(enumName string) []string {
	// ตัวอย่างของการ return ค่าของแต่ละ ENUM type
	switch enumName {
	case "transaction_type":
		return []string{"INCOME", "EXPENSE", "ANY"}
	// case "payment_status":
	// 	return []string{"PENDING", "COMPLETED", "FAILED"}
	// เพิ่ม ENUM type อื่นๆ ที่คุณต้องการที่นี่
	default:
		return []string{}
	}
}
