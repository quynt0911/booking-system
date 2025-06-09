package model

import "time"

type OffTime struct {
    ID        int       `json:"id" db:"id"`
    ExpertID  int       `json:"expert_id" db:"expert_id"`
    StartDate time.Time `json:"start_date" db:"start_date"`
    EndDate   time.Time `json:"end_date" db:"end_date"`
    Reason    string    `json:"reason" db:"reason"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}