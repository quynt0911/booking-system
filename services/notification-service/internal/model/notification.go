package model

import (
	"time"
)

type NotificationType string

const (
	EmailNotification     NotificationType = "email"
	TelegramNotification  NotificationType = "telegram"
	WebSocketNotification NotificationType = "websocket"
)

type Notification struct {
	ID        uint             `json:"id" gorm:"primaryKey"`
	Type      NotificationType `json:"type"`
	Recipient string           `json:"recipient"`
	Subject   string           `json:"subject"`
	Content   string           `json:"content"`
	Status    string           `json:"status"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}
