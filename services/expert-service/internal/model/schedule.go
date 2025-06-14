package model

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ID             uuid.UUID `json:"id"`
	ExpertID       string    `json:"expert_id"`
	UserID         string    `json:"user_id"`
	AvailabilityID string    `json:"availability_id"`
	DayOfWeek      int       `json:"day_of_week"`
	Date           string    `json:"date"`
	StartTime      string    `json:"start_time"`
	EndTime        string    `json:"end_time"`
	Status         string    `json:"status"`
	Title          string    `json:"title"`
	Description    string    `json:"description,omitempty"`
	MeetingLink    string    `json:"meeting_link,omitempty"`
	Notes          string    `json:"notes,omitempty"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
