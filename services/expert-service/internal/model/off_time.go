package model

import (
	"time"

	"github.com/google/uuid"
)

type OffTime struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ExpertID  string    `json:"expert_id" db:"expert_id"`
	StartDate string    `json:"start_date" db:"start_date"`
	EndDate   string    `json:"end_date" db:"end_date"`
	Reason    string    `json:"reason" db:"reason"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
