// Booking model định nghĩa cấu trúc dữ liệu cho booking trong hệ thống tư vấn
package model

import (
	"time"
)

// BookingStatus định nghĩa các trạng thái của booking
type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusConfirmed BookingStatus = "confirmed"
	StatusRejected  BookingStatus = "rejected"
	StatusCancelled BookingStatus = "cancelled"
	StatusCompleted BookingStatus = "completed"
	StatusMissed    BookingStatus = "missed"
)

// BookingType định nghĩa loại hình tư vấn
type BookingType string

const (
	TypeOnline  BookingType = "online"
	TypeOffline BookingType = "offline"
)

// Booking struct định nghĩa cấu trúc của một booking
type Booking struct {
	ID          uint          `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      uint          `json:"user_id" gorm:"not null;index"`
	ExpertID    uint          `json:"expert_id" gorm:"not null;index"`
	StartTime   time.Time     `json:"start_time" gorm:"not null;index"`
	EndTime     time.Time     `json:"end_time" gorm:"not null"`
	Type        BookingType   `json:"type" gorm:"not null;default:'online'"`
	Status      BookingStatus `json:"status" gorm:"not null;default:'pending';index"`
	Notes       string        `json:"notes" gorm:"type:text"`
	Location    string        `json:"location,omitempty"`
	MeetingLink string        `json:"meeting_link,omitempty"`
	CreatedAt   time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	ConfirmedAt *time.Time    `json:"confirmed_at,omitempty"`
	CancelledAt *time.Time    `json:"cancelled_at,omitempty"`
}

// TableName trả về tên bảng trong database
func (Booking) TableName() string {
	return "bookings"
}

// IsActive kiểm tra booking có đang active không
func (b *Booking) IsActive() bool {
	return b.Status == StatusPending || b.Status == StatusConfirmed
}

// CanBeCancelled kiểm tra booking có thể bị hủy không
func (b *Booking) CanBeCancelled() bool {
	if b.Status != StatusPending && b.Status != StatusConfirmed {
		return false
	}
	
	// Chỉ có thể hủy trước 1 tiếng
	oneHourBefore := b.StartTime.Add(-1 * time.Hour)
	return time.Now().Before(oneHourBefore)
}

// CanBeConfirmed kiểm tra booking có thể được xác nhận không
func (b *Booking) CanBeConfirmed() bool {
	return b.Status == StatusPending
}

// IsExpired kiểm tra booking đã hết hạn chưa
func (b *Booking) IsExpired() bool {
	return time.Now().After(b.EndTime) && b.Status == StatusPending
}

// GetDuration trả về thời gian tư vấn (phút)
func (b *Booking) GetDuration() int {
	return int(b.EndTime.Sub(b.StartTime).Minutes())
}

// BookingFilter struct để filter booking
type BookingFilter struct {
	UserID    *uint          `json:"user_id,omitempty"`
	ExpertID  *uint          `json:"expert_id,omitempty"`
	Status    *BookingStatus `json:"status,omitempty"`
	Type      *BookingType   `json:"type,omitempty"`
	StartDate *time.Time     `json:"start_date,omitempty"`
	EndDate   *time.Time     `json:"end_date,omitempty"`
	Page      int            `json:"page" validate:"min=1"`
	Limit     int            `json:"limit" validate:"min=1,max=100"`
	SortBy    string         `json:"sort_by,omitempty"`
	SortOrder string         `json:"sort_order,omitempty"`
}