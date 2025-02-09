package helpers

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WhereConditions(query *gorm.DB, field string, filter interface{}) *gorm.DB {
	switch v := filter.(type) {
	case string:
		if v == "null" {
			return query.Where(field + " IS NULL")
		}
		return query.Where(field, v)
	case []string:
		// ถ้า filter เป็น array
		query = query.Where("1 = 1") // เริ่ม query ด้วยเงื่อนไขพื้นฐาน 1=1 เป็นจริงอยู่แล้ว ต้องการเอา instance เพื่อใช้เงื่อนไขอื่นต่อ เพราะ or ถ้าไม่มีการใช้ where ก่อนจะไม่ทำงาน
		for _, f := range v {
			if f == "null" {
				query = query.Or(field + " IS NULL")
			} else {
				query = query.Or(field, f)
			}
		}
		return query
	default:
		// ค่าอื่นๆ เช่น integer
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
	// เช็คว่า type relations
	switch v := relations.(type) {
	case string:
		query = query.Preload(v)
	case []string:
		for _, relation := range v {
			query = query.Preload(relation)
		}
	default:
		fmt.Println("Invalid type for relations:", v)
	}
	return query
}
