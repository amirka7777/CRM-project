package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserRequest struct {
	Email      string
	Password   string
	Created_at time.Time
}

type UserLoginCheck struct {
	Id           int
	Email        string
	PasswordHash string
}

type Claims struct {
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
}

type CreateContactRequest struct {
	Name  string
	Email string
	Phone string
}
