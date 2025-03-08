package helpers

import (
	"errors"
	"fmt"
	"go_income_outflow/entities"
	"strings"
	"time"

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
		if len(v) > 0 {
			query = query.Where(field+" IN (?)", v)
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
	// เช็ค type relations
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

func DateRange(query *gorm.DB, fromDate, toDate *string, timeIncluded bool, field string) *gorm.DB {
	// 2006-01-02 = YYYY-MM-DD
	// 15:04:05 = HH:mm:ss

	// กำหนดค่า default สำหรับวันที่
	now := time.Now()
	nowFormatted := now.Format("2006-01-02")
	if !timeIncluded {
		// ถ้าไม่ต้องการเวลา ก็จะเป็นแค่วันที่
		nowFormatted = now.Format("2006-01-02")
	} else {
		// ถ้าต้องการเวลา ก็จะมีเวลา
		nowFormatted = now.Format("2006-01-02 15:04:05")
	}

	// ถ้าไม่ได้ระบุวันที่ to_date จะเป็นวันนี้
	if toDate == nil {
		toDate = &nowFormatted
	}

	// ถ้าระบุจากวันที่และจากวันที่นั้นมากกว่าถึงวันที่ จะเปลี่ยนให้เท่ากับถึงวันที่
	if fromDate != nil && *fromDate > *toDate {
		*fromDate = *toDate
	}

	// ใช้เวลาในกรณีที่ต้องการเวลา
	if timeIncluded {
		if fromDate != nil {
			// ถ้ามี from_date ก็ให้กรองจากวันที่และเวลานั้น
			query = query.Where(fmt.Sprintf("%s >= ?", field), *fromDate).Where(fmt.Sprintf("%s <= ?", field), *toDate)
		} else {
			// ถ้าไม่มี from_date ก็ให้กรองแค่ถึงวันที่เท่านั้น
			query = query.Where(fmt.Sprintf("%s <= ?", field), *toDate)
		}
	} else {
		// ถ้าไม่มีเวลา ให้กรองจากวันที่เริ่มต้นถึงวันที่สุดท้ายของวัน
		if fromDate != nil {
			// ถ้ามี from_date จะตั้งเวลาเป็น 00:00:00
			fromDateParsed, _ := time.Parse("2006-01-02", *fromDate)
			fromDateParsed = fromDateParsed.Add(time.Hour * 0) // เวลาเริ่มต้นที่ 00:00:00
			*fromDate = fromDateParsed.Format("2006-01-02 15:04:05")
		}
		// ตั้งเวลาให้ to_date เป็น 23:59:59
		toDateParsed, _ := time.Parse("2006-01-02", *toDate)
		toDateParsed = toDateParsed.Add(time.Hour * 24).Add(time.Nanosecond * -1) // เวลาเป็น 23:59:59
		*toDate = toDateParsed.Format("2006-01-02 15:04:05")

		// กรองช่วงวันที่
		if fromDate != nil {
			query = query.Where(fmt.Sprintf("%s >= ?", field), *fromDate).Where(fmt.Sprintf("%s <= ?", field), *toDate)
		} else {
			query = query.Where(fmt.Sprintf("%s <= ?", field), *toDate)
		}
	}

	return query
}

// แปลงจาก string (YYYY-MM-DD) เป็น time.Time
func ParseDate(dateStr string) (time.Time, error) {
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format: %v", err)
	}
	return parsedDate, nil
}

func UintToPointer(u uint) *uint {
	if u == 0 {
		return nil
	}
	return &u
}

func AuthUser(ctx *gin.Context) (*entities.User, error) {
	// ดึง user จาก context
	user, exists := ctx.Get("user")
	if !exists {
		return nil, errors.New("user not found in context")
	}

	// แปลง user เป็นประเภท *entities.User
	userObj, ok := user.(entities.User)
	if !ok {
		return nil, errors.New("invalid user type in context")
	}

	// คืนค่า pointer ของ user
	return &userObj, nil
}

func ExtractJWT(c *gin.Context) string {
	tokenString := c.GetHeader("Authorization")
	if tokenString != "" {
		if strings.HasPrefix(tokenString, "Bearer ") {
			return strings.TrimPrefix(tokenString, "Bearer ") // ตัด Bearer ออกจาก prefix
		}
		return tokenString
	}

	tokenString, err := c.Cookie("Authorization")
	if err == nil {
		return tokenString
	}

	tokenString = c.Query("token")
	if tokenString != "" {
		return tokenString
	}

	return ""
}

func ExtractRefreshToken(c *gin.Context) string {
	tokenString := c.GetHeader("RefreshToken")
	if tokenString != "" {
		if strings.HasPrefix(tokenString, "Bearer ") {
			return strings.TrimPrefix(tokenString, "Bearer ") // ตัด Bearer ออกจาก prefix
		}
		return tokenString
	}

	tokenString, err := c.Cookie("RefreshToken")
	if err == nil {
		return tokenString
	}

	tokenString = c.Query("refresh_token")
	if tokenString != "" {
		return tokenString
	}

	return ""
}
