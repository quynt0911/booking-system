package model

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
    ID        uuid.UUID `json:"id"`
    ExpertID  int       `json:"expert_id"`
    DayOfWeek int       `json:"day_of_week"`
    StartTime string    `json:"start_time"`
    EndTime   string    `json:"end_time"`
    IsActive  bool      `json:"is_active"`
    CreatedAt time.Time `json:"created_at"`
}