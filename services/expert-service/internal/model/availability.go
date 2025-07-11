package model

import (
	"time"

	"github.com/google/uuid"
)

type Availability struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ExpertID  string    `json:"expert_id" db:"expert_id"`
	Date      string    `json:"date" db:"date"`
	StartTime string    `json:"start_time" db:"start_time"`
	EndTime   string    `json:"end_time" db:"end_time"`
	IsBooked  bool      `json:"is_booked" db:"is_booked"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type AvailabilitySlot struct {
	Date      string `json:"date"`       // "2024-03-25"
	StartTime string `json:"start_time"` // "09:00"
	EndTime   string `json:"end_time"`   // "10:00"
}
