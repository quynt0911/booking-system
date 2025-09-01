package model

type NotificationRequest struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}
