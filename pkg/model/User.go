package model

import "time"

type (
	UserRequest struct {
		Email    string
		Password string
		Name     string
	}
)

type User struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time

	Email    string
	Password string
	Name     string
}

type UserResponse struct {
	ID uint `json:"id"`

	Name  string `json:"name"`
	Email string `json:"email"`
}
