package model

import (
	"time"
)

type Settings struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	UserID           uint      `json:"user_id"`
	EmailEnabled     bool      `json:"email_enabled"`
	TelegramEnabled  bool      `json:"telegram_enabled"`
	WebSocketEnabled bool      `json:"websocket_enabled"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
