package model

import (
	"time"
)

type User struct {
	ID          string     `json:"id" db:"id"`
	Email       string     `json:"email" db:"email" validate:"required,email"`
	Password    string     `json:"-" db:"password" validate:"required,password"`
	FirstName   string     `json:"first_name" db:"first_name" validate:"required,min=2,max=50"`
	LastName    string     `json:"last_name" db:"last_name" validate:"required,min=2,max=50"`
	Role        string     `json:"role" db:"role" validate:"required,oneof=user expert admin"`
	Phone       *string    `json:"phone" db:"phone" validate:"omitempty,phone"`
	AvatarURL   *string    `json:"avatar_url" db:"avatar_url"`
	Bio         *string    `json:"bio" db:"bio" validate:"omitempty,max=500"`
	Gender      *string    `json:"gender" db:"gender" validate:"omitempty,oneof=male female other"`
	DateOfBirth *time.Time `json:"date_of_birth" db:"date_of_birth"`
	IsVerified  bool       `json:"is_verified" db:"is_verified"`
	IsActive    bool       `json:"is_active" db:"is_active"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type UserRole string

const (
	RoleUser   UserRole = "user"
	RoleExpert UserRole = "expert"
	RoleAdmin  UserRole = "admin"
)

// UserProfile represents user profile information
type UserProfile struct {
	ID          string     `json:"id"`
	Email       string     `json:"email"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Role        string     `json:"role"`
	Phone       *string    `json:"phone"`
	AvatarURL   *string    `json:"avatar_url"`
	Bio         *string    `json:"bio"`
	Gender      *string    `json:"gender"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	IsVerified  bool       `json:"is_verified"`
	CreatedAt   time.Time  `json:"created_at"`
}

// ToProfile converts User to UserProfile (removes sensitive data)
func (u *User) ToProfile() *UserProfile {
	return &UserProfile{
		ID:          u.ID,
		Email:       u.Email,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Role:        u.Role,
		Phone:       u.Phone,
		AvatarURL:   u.AvatarURL,
		Bio:         u.Bio,
		Gender:      u.Gender,
		DateOfBirth: u.DateOfBirth,
		IsVerified:  u.IsVerified,
		CreatedAt:   u.CreatedAt,
	}
}

// GetFullName returns user's full name
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// IsExpert checks if user is an expert
func (u *User) IsExpert() bool {
	return u.Role == string(RoleExpert)
}

// IsAdmin checks if user is an admin
func (u *User) IsAdmin() bool {
	return u.Role == string(RoleAdmin)
}

// CanManageBookings checks if user can manage bookings
func (u *User) CanManageBookings() bool {
	return u.IsExpert() || u.IsAdmin()
}

