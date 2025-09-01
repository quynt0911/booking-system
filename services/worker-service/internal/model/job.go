package model

import (
	"time"
)

type JobType string

const (
	ReminderJob     JobType = "reminder"
	CleanupJob      JobType = "cleanup"
	EmailJob        JobType = "email"
	StatusUpdateJob JobType = "status_update"
)

type Job struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Type      JobType   `json:"type"`
	Payload   string    `json:"payload"`
	Status    string    `json:"status"`
	Schedule  string    `json:"schedule"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
