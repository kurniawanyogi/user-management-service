package model

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type User struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"-" db:"password"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	Email     string    `json:"email" db:"email"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

type RegistrationUserRequest struct {
	Username  string  `json:"username" validate:"required,min=3,max=32"`
	Password  string  `json:"password" validate:"required,min=8,max=32"`
	FirstName string  `json:"firstName" validate:"required,min=1,max=100"`
	LastName  *string `json:"lastName" validate:"omitempty,max=100"`
	Email     string  `json:"email" validate:"required,email,max=100"`
}

type UpdateUserRequest struct {
	Id        int64
	FirstName string  `json:"firstName" validate:"required,min=1,max=100"`
	LastName  *string `json:"lastName" validate:"omitempty,max=100"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
