// Booking model định nghĩa cấu trúc dữ liệu cho booking trong hệ thống tư vấn
package model

import (
	"time"

	"github.com/google/uuid"
)

// BookingType represents the type of booking
type BookingType string

const (
	TypeOnline  BookingType = "online"
	TypeOffline BookingType = "offline"
)

// BookingStatus represents the status of a booking
type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusRejected  BookingStatus = "rejected"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusCompleted BookingStatus = "completed"
	BookingStatusMissed    BookingStatus = "missed"
)

// Booking represents a booking record
type Booking struct {
	ID              uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID          uuid.UUID     `json:"user_id" gorm:"type:uuid"`
	ExpertID        uuid.UUID     `json:"expert_id" gorm:"type:uuid"`
	ScheduledTime   time.Time     `json:"scheduled_datetime" gorm:"column:scheduled_datetime"`
	DurationMinutes int           `json:"duration_minutes" gorm:"default:60"`
	MeetingType     BookingType   `json:"meeting_type" gorm:"default:'online'"`
	MeetingURL      string        `json:"meeting_url,omitempty"`
	MeetingAddress  string        `json:"meeting_address,omitempty"`
	Notes           string        `json:"notes,omitempty"`
	Status          BookingStatus `json:"status" gorm:"default:'pending'"`
	Price           float64       `json:"price,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	ConfirmedAt     *time.Time    `json:"confirmed_at,omitempty"`
	CancelledAt     *time.Time    `json:"cancelled_at,omitempty"`
	CompletedAt     *time.Time    `json:"completed_at,omitempty"`
}

// IsValidBookingStatus checks if the given status is valid
func IsValidBookingStatus(status string) bool {
	switch BookingStatus(status) {
	case BookingStatusPending, BookingStatusConfirmed, BookingStatusRejected,
		BookingStatusCancelled, BookingStatusCompleted, BookingStatusMissed:
		return true
	default:
		return false
	}
}

// TableName returns the table name in the database
func (Booking) TableName() string {
	return "bookings"
}

// IsActive checks if the booking is currently active
func (b *Booking) IsActive() bool {
	return b.Status == BookingStatusPending || b.Status == BookingStatusConfirmed
}

// CanBeCancelled checks if the booking can be cancelled
func (b *Booking) CanBeCancelled() bool {
	if b.Status != BookingStatusPending && b.Status != BookingStatusConfirmed {
		return false
	}

	// Can only cancel 1 hour before the scheduled time
	oneHourBefore := b.ScheduledTime.Add(-1 * time.Hour)
	return time.Now().Before(oneHourBefore)
}

// CanBeConfirmed checks if the booking can be confirmed
func (b *Booking) CanBeConfirmed() bool {
	return b.Status == BookingStatusPending
}

// IsExpired checks if the booking has expired
func (b *Booking) IsExpired() bool {
	endTime := b.ScheduledTime.Add(time.Duration(b.DurationMinutes) * time.Minute)
	return time.Now().After(endTime) && b.Status == BookingStatusPending
}

// GetEndTime returns the end time of the booking
func (b *Booking) GetEndTime() time.Time {
	return b.ScheduledTime.Add(time.Duration(b.DurationMinutes) * time.Minute)
}

// BookingFilter struct for filtering bookings
type BookingFilter struct {
	UserID    *uuid.UUID     `json:"user_id,omitempty"`
	ExpertID  *uuid.UUID     `json:"expert_id,omitempty"`
	Status    *BookingStatus `json:"status,omitempty"`
	Type      *BookingType   `json:"type,omitempty"`
	StartDate *time.Time     `json:"start_date,omitempty"`
	EndDate   *time.Time     `json:"end_date,omitempty"`
	Page      int            `json:"page" validate:"min=1"`
	Limit     int            `json:"limit" validate:"min=1,max=100"`
	SortBy    string         `json:"sort_by,omitempty"`
	SortOrder string         `json:"sort_order,omitempty"`
}
