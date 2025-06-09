package model

import (
	"time"
)
type Expert struct {
    ID             int       `json:"id"`
    Name           string    `json:"name"`
    Email          string    `json:"email"`
    Specialization string    `json:"specialization"`
    Status         string    `json:"status"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}