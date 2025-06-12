package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Expert struct {
	ID              uuid.UUID      `json:"id" db:"id"`
	UserID          uuid.UUID      `json:"user_id" db:"user_id"`
	Specialization  string         `json:"specialization" db:"specialization"`
	ExperienceYears int            `json:"experience_years" db:"experience_years"`
	HourlyRate      float64        `json:"hourly_rate" db:"hourly_rate"`
	Certifications  pq.StringArray `json:"certifications" db:"certifications"`
	IsAvailable     bool           `json:"is_available" db:"is_available"`
	Rating          float64        `json:"rating" db:"rating"`
	TotalReviews    int            `json:"total_reviews" db:"total_reviews"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
}
