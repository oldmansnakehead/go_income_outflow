package common

import "time"

type Model struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type CommonRequest struct {
	With []string `json:"with" query:"with"` // รับจากทั้ง body และ query
}
