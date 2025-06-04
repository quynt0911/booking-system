package model

import "time"

type Expert struct {
    ID             int       `json:"id" db:"id"`
    Name           string    `json:"name" db:"name"`
    Email          string    `json:"email" db:"email"`
    Specialization string    `json:"specialization" db:"specialization"`
    Status         string    `json:"status" db:"status"`
    CreatedAt      time.Time `json:"created_at" db:"created_at"`
    UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}