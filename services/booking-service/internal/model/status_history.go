// Booking model định nghĩa cấu trúc dữ liệu cho booking trong hệ thống tư vấn
package model

import (
	"time"

	"github.com/google/uuid"
)

// StatusHistory represents a record of booking status changes
type StatusHistory struct {
	ID        uuid.UUID     `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	BookingID uuid.UUID     `json:"booking_id" gorm:"type:uuid"`
	Status    BookingStatus `json:"status"`
	ChangedBy uuid.UUID     `json:"changed_by" gorm:"type:uuid"`
	ChangedAt time.Time     `json:"changed_at"`
	Note      string        `json:"note,omitempty"`
}

// TableName trả về tên bảng trong database
func (StatusHistory) TableName() string {
	return "status_histories"
}

// ChangeType constants
const (
	ChangeTypeUser   = "user"
	ChangeTypeExpert = "expert"
	ChangeTypeSystem = "system"
	ChangeTypeAdmin  = "admin"
)

// StatusTransition struct định nghĩa các chuyển đổi trạng thái hợp lệ
type StatusTransition struct {
	From      BookingStatus
	To        BookingStatus
	AllowedBy []string // user, expert, system, admin
}

// ValidTransitions định nghĩa các chuyển đổi trạng thái hợp lệ
var ValidTransitions = []StatusTransition{
	// Từ pending
	{BookingStatusPending, BookingStatusConfirmed, []string{ChangeTypeExpert, ChangeTypeAdmin}},
	{BookingStatusPending, BookingStatusRejected, []string{ChangeTypeExpert, ChangeTypeAdmin}},
	{BookingStatusPending, BookingStatusCancelled, []string{ChangeTypeUser, ChangeTypeExpert, ChangeTypeAdmin}},

	// Từ confirmed
	{BookingStatusConfirmed, BookingStatusCompleted, []string{ChangeTypeExpert, ChangeTypeAdmin}},
	{BookingStatusConfirmed, BookingStatusCancelled, []string{ChangeTypeUser, ChangeTypeExpert, ChangeTypeAdmin}},

	// Từ rejected (chỉ admin có thể thay đổi)
	{BookingStatusRejected, BookingStatusPending, []string{ChangeTypeAdmin}},
	{BookingStatusRejected, BookingStatusConfirmed, []string{ChangeTypeAdmin}},

	// Từ cancelled (chỉ admin có thể thay đổi)
	{BookingStatusCancelled, BookingStatusPending, []string{ChangeTypeAdmin}},
	{BookingStatusCancelled, BookingStatusConfirmed, []string{ChangeTypeAdmin}},
}

// IsValidTransition kiểm tra chuyển đổi trạng thái có hợp lệ không
func IsValidTransition(from, to BookingStatus, changeType string) bool {
	for _, transition := range ValidTransitions {
		if transition.From == from && transition.To == to {
			for _, allowedType := range transition.AllowedBy {
				if allowedType == changeType {
					return true
				}
			}
		}
	}
	return false
}

// GetValidNextStatuses trả về danh sách trạng thái có thể chuyển đến
func GetValidNextStatuses(from BookingStatus, changeType string) []BookingStatus {
	var validStatuses []BookingStatus

	for _, transition := range ValidTransitions {
		if transition.From == from {
			for _, allowedType := range transition.AllowedBy {
				if allowedType == changeType {
					validStatuses = append(validStatuses, transition.To)
					break
				}
			}
		}
	}

	return validStatuses
}

// StatusHistoryFilter struct để filter lịch sử trạng thái
type StatusHistoryFilter struct {
	BookingID *uuid.UUID     `json:"booking_id,omitempty"`
	ChangedBy *uuid.UUID     `json:"changed_by,omitempty"`
	Status    *BookingStatus `json:"status,omitempty"`
	StartDate *time.Time     `json:"start_date,omitempty"`
	EndDate   *time.Time     `json:"end_date,omitempty"`
	Page      int            `json:"page" validate:"min=1"`
	Limit     int            `json:"limit" validate:"min=1,max=100"`
}
