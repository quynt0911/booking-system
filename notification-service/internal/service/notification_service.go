package service

import "notification-service/internal/model"

func FetchRecentNotifications() []model.Notification {
	return []model.Notification{
		{UserID: "123", Message: "Your appointment is confirmed", Read: false},
	}
}