package model

import "time"

type User struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	Email    string
	Password string
	Name     string
}
