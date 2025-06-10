package model

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleExpert UserRole = "expert"
	RoleUser   UserRole = "user"
)

type User struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email         string    `json:"email" gorm:"unique;not null"`
	PasswordHash  string    `json:"-" gorm:"column:password_hash;not null"`
	FullName      string    `json:"fullname" gorm:"column:full_name;not null"`
	Phone         string    `json:"phone"`
	Image         string    `json:"image"`
	Gender        string    `json:"gender"`
	Description   string    `json:"description"`
	Role          UserRole  `json:"role" gorm:"default:'user'"`
	EmailVerified bool      `json:"email_verified" gorm:"default:false"`
	CreatedAt     time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	FullName    string `json:"fullname" binding:"required"`
	Phone       string `json:"phone"`
	Gender      string `json:"gender"`
	Role        string `json:"role" binding:"required,oneof=admin expert user"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	FullName    string `json:"fullname"`
	Phone       string `json:"phone"`
	Image       string `json:"image"`
	Gender      string `json:"gender"`
	Description string `json:"description"`
}
