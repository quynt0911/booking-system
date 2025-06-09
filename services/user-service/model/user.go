package model

type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleExpert UserRole = "expert"
	RoleUser   UserRole = "user"
)

type User struct {
	ID            string   `json:"id" gorm:"primaryKey"`
	Email         string   `json:"email" gorm:"unique;not null"`
	PasswordHash  string   `json:"-" gorm:"column:password_hash;not null"`
	FullName      string   `json:"full_name" gorm:"column:fullname;not null"`
	Role          UserRole `json:"role" gorm:"type:varchar(20);default:'user'"`
	Phone         string   `json:"phone" gorm:"type:varchar(20)"`
	EmailVerified bool     `json:"email_verified" gorm:"default:false"`
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
