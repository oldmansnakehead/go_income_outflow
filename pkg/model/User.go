package model

import "time"

type (
	User struct {
		ID        uint
		CreatedAt time.Time
		UpdatedAt time.Time

		Email    string
		Password string
		Name     string
	}

	UserLoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	UserRequest struct {
		Email    string `json:"email" validate:"required,email,uniqueEmail"`
		Password string `json:"password" validate:"required,min=8"`
		Name     string `json:"name" validate:"required"`
	}

	UserResponse struct {
		ID uint `json:"id"`

		Name  string `json:"name"`
		Email string `json:"email"`
	}
)
