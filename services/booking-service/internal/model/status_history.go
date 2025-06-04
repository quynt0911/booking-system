// Booking model định nghĩa cấu trúc dữ liệu cho booking trong hệ thống tư vấn
package model

import (
	"time"
)

// StatusHistory struct lưu lịch sử thay đổi trạng thái của booking
type StatusHistory struct {
	ID          uint          `json:"id" gorm:"primaryKey;autoIncrement"`
	BookingID   uint          `json:"booking_id" gorm:"not null;index"`
	OldStatus   BookingStatus `json:"old_status" gorm:"not null"`
	NewStatus   BookingStatus `json:"new_status" gorm:"not null"`
	ChangedBy   uint          `json:"changed_by" gorm:"not null"` // UserID của người thay đổi
	ChangeType  string        `json:"change_type" gorm:"not null"` // user, expert, system, admin
	Reason      string        `json:"reason,omitempty" gorm:"type:text"`
	Notes       string        `json:"notes,omitempty" gorm:"type:text"`
	CreatedAt   time.Time     `json:"created_at" gorm:"autoCreateTime"`
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
	From    BookingStatus
	To      BookingStatus
	AllowedBy []string // user, expert, system, admin
}

// ValidTransitions định nghĩa các chuyển đổi trạng thái hợp lệ
var ValidTransitions = []StatusTransition{
	// Từ pending
	{StatusPending, StatusConfirmed, []string{ChangeTypeExpert, ChangeTypeAdmin}},
	{StatusPending, StatusRejected, []string{ChangeTypeExpert, ChangeTypeAdmin}},
	{StatusPending, StatusCancelled, []string{ChangeTypeUser, ChangeTypeExpert, ChangeTypeAdmin}},
	{StatusPending, StatusMissed, []string{ChangeTypeSystem, ChangeTypeAdmin}},
	
	// Từ confirmed
	{StatusConfirmed, StatusCompleted, []string{ChangeTypeExpert, ChangeTypeAdmin}},
	{StatusConfirmed, StatusCancelled, []string{ChangeTypeUser, ChangeTypeExpert, ChangeTypeAdmin}},
	{StatusConfirmed, StatusMissed, []string{ChangeTypeSystem, ChangeTypeAdmin}},
	
	// Từ rejected (chỉ admin có thể thay đổi)
	{StatusRejected, StatusPending, []string{ChangeTypeAdmin}},
	{StatusRejected, StatusConfirmed, []string{ChangeTypeAdmin}},
	
	// Từ cancelled (chỉ admin có thể thay đổi)
	{StatusCancelled, StatusPending, []string{ChangeTypeAdmin}},
	{StatusCancelled, StatusConfirmed, []string{ChangeTypeAdmin}},
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
	BookingID  *uint   `json:"booking_id,omitempty"`
	ChangedBy  *uint   `json:"changed_by,omitempty"`
	ChangeType *string `json:"change_type,omitempty"`
	OldStatus  *BookingStatus `json:"old_status,omitempty"`
	NewStatus  *BookingStatus `json:"new_status,omitempty"`
	StartDate  *time.Time `json:"start_date,omitempty"`
	EndDate    *time.Time `json:"end_date,omitempty"`
	Page       int     `json:"page" validate:"min=1"`
	Limit      int     `json:"limit" validate:"min=1,max=100"`
}