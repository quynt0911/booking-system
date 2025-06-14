package model

import (
	"time"

	"github.com/google/uuid"
)

type OffTime struct {
	ID            uuid.UUID `json:"id" db:"id"`
	ExpertID      uuid.UUID `json:"expert_id" db:"expert_id"`
	StartDateTime time.Time `json:"start_datetime" db:"start_datetime"`
	EndDateTime   time.Time `json:"end_datetime" db:"end_datetime"`
	Reason        string    `json:"reason" db:"reason"`
	IsRecurring   bool      `json:"is_recurring" db:"is_recurring"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}
