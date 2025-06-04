package model

import "time"

type Schedule struct {
    ID        int       `json:"id" db:"id"`
    ExpertID  int       `json:"expert_id" db:"expert_id"`
    DayOfWeek int       `json:"day_of_week" db:"day_of_week"`
    StartTime string    `json:"start_time" db:"start_time"`
    EndTime   string    `json:"end_time" db:"end_time"`
    IsActive  bool      `json:"is_active" db:"is_active"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}