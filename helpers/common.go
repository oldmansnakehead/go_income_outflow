package helpers

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SingleOrMultiple(query *gorm.DB, field string, filter interface{}) *gorm.DB {
	switch v := filter.(type) {
	case string:
		if v == "null" {
			return query.Where(field + " IS NULL")
		}
		return query.Where(field, v)
	case []string:
		// ถ้า filter เป็น array
		query = query.Where("1 = 1") // เริ่ม query ด้วยเงื่อนไขพื้นฐาน
		for _, f := range v {
			if f == "null" {
				query = query.Or(field + " IS NULL")
			} else {
				query = query.Or(field, f)
			}
		}
		return query
	default:
		// ค่าอื่นๆ เช่น integer สามารถขยายได้ตามความต้องการ
		return query.Where(field, v)
	}
}

// GetFiltersFromQuery รับค่า filters เช่น employee_id, user_id, with จาก query string
// output = { "user_id": ["2", "1"], "with": ["User", "Employee"] }
func ParseQueryString(ctx *gin.Context) map[string]interface{} {
	filters := make(map[string]interface{})

	// ดึงข้อมูลทั้งหมดจาก query string
	queryParams := ctx.Request.URL.Query()

	// วนลูปเพื่อจัดการกับ query string ที่ส่งมาทุกตัว
	for key, values := range queryParams {
		// ลบ '[]' ออกจาก key ถ้ามี
		key = strings.TrimSuffix(key, "[]")

		// ถ้าค่ามีหลายค่า (array)
		if len(values) > 1 {
			filters[key] = values
		} else if len(values) == 1 {
			// ถ้ามีแค่ตัวเดียว
			filters[key] = values[0]
		}
	}

	return filters
}

func WithRelations(query *gorm.DB, relations interface{}) *gorm.DB {
	// ตรวจสอบ type ของ relations
	fmt.Printf("Type of relations: %T\n", relations) // พิมพ์ type ของ relations

	// เช็คว่า relations เป็น string หรือไม่
	switch v := relations.(type) {
	case string:
		// ถ้าเป็น string, ให้ preload ค่านั้น
		fmt.Println("Preloading relation:", v)
		query = query.Preload(v)
	case []string:
		// ถ้าเป็น []string, วนลูปเพื่อ preload แต่ละ relation
		for _, relation := range v {
			fmt.Println("Preloading relation:", relation)
			query = query.Preload(relation)
		}
	default:
		// หากไม่ใช่ string หรือ []string, ให้ไม่ทำอะไร
		fmt.Println("Invalid type for relations:", v)
	}
	return query
}
