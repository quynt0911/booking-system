package service

import (
	"context"
	"fmt"
	"notification-service/internal/config"
	"notification-service/internal/model"
	"notification-service/internal/repository"
)

type NotificationService struct {
	repo             *repository.NotificationRepository
	emailService     *EmailService
	telegramService  *TelegramService
	websocketService *WebSocketService
}

func NewNotificationService(cfg *config.Config) *NotificationService {
	return &NotificationService{
		repo:             repository.NewNotificationRepository(cfg.DB),
		emailService:     NewEmailService(cfg),
		telegramService:  NewTelegramService(cfg),
		websocketService: NewWebSocketService(nil), // Hub will be set later
	}
}

func (s *NotificationService) SendNotification(ctx context.Context, notification *model.Notification) error {
	// Save notification to database
	if err := s.repo.Create(ctx, notification); err != nil {
		return err
	}

	// Send notification based on type
	switch notification.Type {
	case model.EmailNotification:
		return s.emailService.Send(notification.Recipient, notification.Subject, notification.Content)
	case model.TelegramNotification:
		return s.telegramService.Send(notification.Recipient, notification.Content)
	case model.WebSocketNotification:
		return s.websocketService.Send(notification.Recipient, notification.Content)
	default:
		return fmt.Errorf("unsupported notification type: %s", notification.Type)
	}
}

func FetchRecentNotifications() []model.Notification {
	return []model.Notification{
		{UserID: "123", Message: "Your appointment is confirmed", Read: false},
	}
}
