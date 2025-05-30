package model

type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleExpert UserRole = "expert"
	RoleUser   UserRole = "user"
)

type User struct {
	ID       string   `json:"id" gorm:"primaryKey"`
	Email    string   `json:"email" gorm:"unique;not null"`
	Password string   `json:"-"`
	FullName string   `json:"full_name"`
	Role     UserRole `json:"role"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=admin expert user"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	FullName string `json:"full_name" binding:"required"`
}
