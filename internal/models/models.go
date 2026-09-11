package models

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	PasswordHash string    `json:"-" gorm:"column:password_hash"`
}

type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}
