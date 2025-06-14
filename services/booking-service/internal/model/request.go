// Booking model định nghĩa cấu trúc dữ liệu cho booking trong hệ thống tư vấn
package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CreateBookingRequest struct for creating a new booking
type CreateBookingRequest struct {
	ExpertID        uuid.UUID `json:"expert_id" binding:"required"`
	ScheduledTime   time.Time `json:"start_time" binding:"required"`
	EndTime         time.Time `json:"end_time" binding:"required"`
	DurationMinutes int       `json:"duration_minutes"`
	MeetingType     string    `json:"meeting_type"`
	MeetingURL      *string   `json:"meeting_url,omitempty"`
	MeetingAddress  *string   `json:"meeting_address,omitempty"`
	Notes           string    `json:"notes"`
	Status          string    `json:"status"`
}

// Validate validates the create booking request
func (req *CreateBookingRequest) Validate() error {
	if req.DurationMinutes < 15 || req.DurationMinutes > 480 {
		return fmt.Errorf("duration must be between 15 and 480 minutes")
	}
	if req.MeetingType != string(TypeOnline) && req.MeetingType != string(TypeOffline) {
		return fmt.Errorf("invalid booking type")
	}
	if req.MeetingType == string(TypeOffline) && req.MeetingAddress == nil {
		return fmt.Errorf("meeting address is required for offline booking")
	}
	if req.MeetingType == string(TypeOnline) && req.MeetingURL == nil {
		return fmt.Errorf("meeting URL is required for online booking")
	}
	return nil
}

// UpdateBookingRequest struct for updating a booking
type UpdateBookingRequest struct {
	ScheduledTime   *time.Time   `json:"scheduled_datetime,omitempty" validate:"omitempty"`
	DurationMinutes *int         `json:"duration_minutes,omitempty" validate:"omitempty,min=15,max=480"`
	MeetingType     *BookingType `json:"meeting_type,omitempty" validate:"omitempty,oneof=online offline"`
	Notes           *string      `json:"notes,omitempty" validate:"omitempty,max=1000"`
	MeetingAddress  *string      `json:"meeting_address,omitempty" validate:"omitempty,max=255"`
	MeetingURL      *string      `json:"meeting_url,omitempty" validate:"omitempty,max=255,url"`
	Price           *float64     `json:"price,omitempty" validate:"omitempty,min=0"`
}

// UpdateBookingStatusRequest struct for changing booking status
type UpdateBookingStatusRequest struct {
	Status BookingStatus `json:"status" validate:"required,oneof=pending confirmed rejected cancelled completed missed"`
	Reason string        `json:"reason,omitempty" validate:"max=500"`
	Notes  string        `json:"notes,omitempty" validate:"max=1000"`
}

// CancelBookingRequest struct for cancelling a booking
type CancelBookingRequest struct {
	Reason string `json:"reason" validate:"required,max=500"`
	Notes  string `json:"notes,omitempty" validate:"max=1000"`
}

// ConfirmBookingRequest struct for confirming a booking
type ConfirmBookingRequest struct {
	MeetingURL string `json:"meeting_url,omitempty" validate:"omitempty,max=255,url"`
	Notes      string `json:"notes,omitempty" validate:"max=1000"`
}

// RejectBookingRequest struct for rejecting a booking
type RejectBookingRequest struct {
	Reason string `json:"reason" validate:"required,max=500"`
	Notes  string `json:"notes,omitempty" validate:"max=1000"`
}

// CompleteBookingRequest struct for completing a booking
type CompleteBookingRequest struct {
	Notes string `json:"notes,omitempty" validate:"max=1000"`
}

// GetBookingsRequest struct for getting a list of bookings
type GetBookingsRequest struct {
	UserID    *uuid.UUID     `json:"user_id,omitempty" query:"user_id"`
	ExpertID  *uuid.UUID     `json:"expert_id,omitempty" query:"expert_id"`
	Status    *BookingStatus `json:"status,omitempty" query:"status" validate:"omitempty,oneof=pending confirmed rejected cancelled completed missed"`
	Type      *BookingType   `json:"type,omitempty" query:"type" validate:"omitempty,oneof=online offline"`
	StartDate *time.Time     `json:"start_date,omitempty" query:"start_date"`
	EndDate   *time.Time     `json:"end_date,omitempty" query:"end_date"`
	Page      int            `json:"page" query:"page" validate:"min=1"`
	Limit     int            `json:"limit" query:"limit" validate:"min=1,max=100"`
	SortBy    string         `json:"sort_by,omitempty" query:"sort_by" validate:"omitempty,oneof=created_at scheduled_datetime status"`
	SortOrder string         `json:"sort_order,omitempty" query:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// GetHistoryRequest struct for getting booking history
type GetHistoryRequest struct {
	BookingID  *uuid.UUID     `json:"booking_id,omitempty" query:"booking_id"`
	ChangedBy  *uuid.UUID     `json:"changed_by,omitempty" query:"changed_by"`
	ChangeType *string        `json:"change_type,omitempty" query:"change_type" validate:"omitempty,oneof=user expert system admin"`
	OldStatus  *BookingStatus `json:"old_status,omitempty" query:"old_status"`
	NewStatus  *BookingStatus `json:"new_status,omitempty" query:"new_status"`
	StartDate  *time.Time     `json:"start_date,omitempty" query:"start_date"`
	EndDate    *time.Time     `json:"end_date,omitempty" query:"end_date"`
	Page       int            `json:"page" query:"page" validate:"min=1"`
	Limit      int            `json:"limit" query:"limit" validate:"min=1,max=100"`
}

// CheckConflictRequest struct for checking time conflicts
type CheckConflictRequest struct {
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	ExpertID  uuid.UUID  `json:"expert_id" validate:"required"`
	StartTime time.Time  `json:"start_time" validate:"required"`
	EndTime   time.Time  `json:"end_time" validate:"required"`
	ExcludeID *uuid.UUID `json:"exclude_id,omitempty"` // Exclude this booking when checking
}

// BookingResponse struct for booking response
type BookingResponse struct {
	*Booking
	ExpertName     string `json:"expert_name,omitempty"`
	UserName       string `json:"user_name,omitempty"`
	CanBeCancelled bool   `json:"can_be_cancelled"`
	CanBeConfirmed bool   `json:"can_be_confirmed"`
	IsExpired      bool   `json:"is_expired"`
	Duration       int    `json:"duration"` // minutes
}

// BookingListResponse struct for booking list response
type BookingListResponse struct {
	Bookings   []BookingResponse `json:"bookings"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// StatusHistoryResponse struct for status history response
type StatusHistoryResponse struct {
	*StatusHistory
	ChangedByName string `json:"changed_by_name,omitempty"`
}

// ConflictCheckResponse struct for conflict check response
type ConflictCheckResponse struct {
	HasConflict      bool              `json:"has_conflict"`
	ConflictBookings []BookingResponse `json:"conflict_bookings,omitempty"`
	Message          string            `json:"message,omitempty"`
}

// BookingStatsResponse struct for booking statistics response
type BookingStatsResponse struct {
	TotalBookings   int64            `json:"total_bookings"`
	StatusBreakdown map[string]int64 `json:"status_breakdown"`
	TypeBreakdown   map[string]int64 `json:"type_breakdown"`
}
