// Booking model định nghĩa cấu trúc dữ liệu cho booking trong hệ thống tư vấn  
package model

import (
	"time"
)

// CreateBookingRequest struct cho request tạo booking mới
type CreateBookingRequest struct {
	ExpertID    uint        `json:"expert_id" validate:"required,min=1"`
	StartTime   time.Time   `json:"start_time" validate:"required"`
	EndTime     time.Time   `json:"end_time" validate:"required"`
	Type        BookingType `json:"type" validate:"required,oneof=online offline"`
	Notes       string      `json:"notes" validate:"max=1000"`
	Location    string      `json:"location,omitempty" validate:"max=255"`
	MeetingLink string      `json:"meeting_link,omitempty" validate:"max=255,url"`
}

// UpdateBookingRequest struct cho request cập nhật booking
type UpdateBookingRequest struct {
	Notes       *string `json:"notes,omitempty" validate:"omitempty,max=1000"`
	Location    *string `json:"location,omitempty" validate:"omitempty,max=255"`
	MeetingLink *string `json:"meeting_link,omitempty" validate:"omitempty,max=255,url"`
}

// UpdateBookingStatusRequest struct cho request thay đổi trạng thái
type UpdateBookingStatusRequest struct {
	Status BookingStatus `json:"status" validate:"required,oneof=pending confirmed rejected cancelled completed missed"`
	Reason string        `json:"reason,omitempty" validate:"max=500"`
	Notes  string        `json:"notes,omitempty" validate:"max=1000"`
}

// CancelBookingRequest struct cho request hủy booking
type CancelBookingRequest struct {
	Reason string `json:"reason" validate:"required,max=500"`
	Notes  string `json:"notes,omitempty" validate:"max=1000"`
}

// ConfirmBookingRequest struct cho request xác nhận booking
type ConfirmBookingRequest struct {
	MeetingLink string `json:"meeting_link,omitempty" validate:"omitempty,max=255,url"`
	Notes       string `json:"notes,omitempty" validate:"max=1000"`
}

// RejectBookingRequest struct cho request từ chối booking
type RejectBookingRequest struct {
	Reason string `json:"reason" validate:"required,max=500"`
	Notes  string `json:"notes,omitempty" validate:"max=1000"`
}

// CompleteBookingRequest struct cho request hoàn thành booking
type CompleteBookingRequest struct {
	Notes string `json:"notes,omitempty" validate:"max=1000"`
}

// GetBookingsRequest struct cho request lấy danh sách booking
type GetBookingsRequest struct {
	UserID    *uint          `json:"user_id,omitempty" query:"user_id"`
	ExpertID  *uint          `json:"expert_id,omitempty" query:"expert_id"`
	Status    *BookingStatus `json:"status,omitempty" query:"status" validate:"omitempty,oneof=pending confirmed rejected cancelled completed missed"`
	Type      *BookingType   `json:"type,omitempty" query:"type" validate:"omitempty,oneof=online offline"`
	StartDate *time.Time     `json:"start_date,omitempty" query:"start_date"`
	EndDate   *time.Time     `json:"end_date,omitempty" query:"end_date"`
	Page      int            `json:"page" query:"page" validate:"min=1"`
	Limit     int            `json:"limit" query:"limit" validate:"min=1,max=100"`
	SortBy    string         `json:"sort_by,omitempty" query:"sort_by" validate:"omitempty,oneof=created_at start_time status"`
	SortOrder string         `json:"sort_order,omitempty" query:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// GetHistoryRequest struct cho request lấy lịch sử booking
type GetHistoryRequest struct {
	BookingID  *uint          `json:"booking_id,omitempty" query:"booking_id"`
	ChangedBy  *uint          `json:"changed_by,omitempty" query:"changed_by"`
	ChangeType *string        `json:"change_type,omitempty" query:"change_type" validate:"omitempty,oneof=user expert system admin"`
	OldStatus  *BookingStatus `json:"old_status,omitempty" query:"old_status"`
	NewStatus  *BookingStatus `json:"new_status,omitempty" query:"new_status"`
	StartDate  *time.Time     `json:"start_date,omitempty" query:"start_date"`
	EndDate    *time.Time     `json:"end_date,omitempty" query:"end_date"`
	Page       int            `json:"page" query:"page" validate:"min=1"`
	Limit      int            `json:"limit" query:"limit" validate:"min=1,max=100"`
}

// CheckConflictRequest struct cho request kiểm tra xung đột thời gian
type CheckConflictRequest struct {
	UserID    *uint     `json:"user_id,omitempty"`
	ExpertID  uint      `json:"expert_id" validate:"required,min=1"`
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required"`
	ExcludeID *uint     `json:"exclude_id,omitempty"` // Loại trừ booking này khi kiểm tra
}

// BookingResponse struct cho response booking
type BookingResponse struct {
	*Booking
	ExpertName     string `json:"expert_name,omitempty"`
	UserName       string `json:"user_name,omitempty"`
	CanBeCancelled bool   `json:"can_be_cancelled"`
	CanBeConfirmed bool   `json:"can_be_confirmed"`
	IsExpired      bool   `json:"is_expired"`
	Duration       int    `json:"duration"` // phút
}

// BookingListResponse struct cho response danh sách booking
type BookingListResponse struct {
	Bookings   []BookingResponse `json:"bookings"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// StatusHistoryResponse struct cho response lịch sử trạng thái
type StatusHistoryResponse struct {
	*StatusHistory
	ChangedByName string `json:"changed_by_name,omitempty"`
}

// ConflictCheckResponse struct cho response kiểm tra xung đột
type ConflictCheckResponse struct {
	HasConflict      bool                `json:"has_conflict"`
	ConflictBookings []BookingResponse   `json:"conflict_bookings,omitempty"`
	Message          string              `json:"message,omitempty"`
}

// BookingStatsResponse struct cho response thống kê booking
type BookingStatsResponse struct {
	TotalBookings     int64             `json:"total_bookings"`
	PendingBookings   int64             `json:"pending_bookings"`
	ConfirmedBookings int64             `json:"confirmed_bookings"`
	CompletedBookings int64             `json:"completed_bookings"`
	CancelledBookings int64             `json:"cancelled_bookings"`
	StatusBreakdown   map[string]int64  `json:"status_breakdown"`
	TypeBreakdown     map[string]int64  `json:"type_breakdown"`
}