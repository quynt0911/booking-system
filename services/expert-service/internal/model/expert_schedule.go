package model

import (
	"time"

	"github.com/google/uuid"
)

type ExpertSchedule struct {
	ID        uuid.UUID `json:"id"`
	ExpertID  uuid.UUID `json:"expert_id"`
	DayOfWeek int       `json:"day_of_week"` // 0=Chủ nhật, 1-6=Thứ 2-7
	StartTime string    `json:"start_time"`  // "09:00"
	EndTime   string    `json:"end_time"`    // "17:00"
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
