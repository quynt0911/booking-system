package notification_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/booking-system/services/notification-service/internal/model"
	"github.com/your-org/booking-system/services/notification-service/internal/service"
)

func TestNotificationService_SendNotification(t *testing.T) {
    // Setup
    mockRepo := newMockNotificationRepository()
    mockEmailService := newMockEmailService()
    mockTelegramService := newMockTelegramService()
    mockWebSocketService := newMockWebSocketService()

    svc := service.NewNotificationService(mockRepo, mockEmailService, mockTelegramService, mockWebSocketService)

    tests := []struct {
        name         string
        notification *model.Notification
        wantErr      bool
    }{
        {
            name: "successful email notification",
            notification: &model.Notification{
                Type:      model.EmailNotification,
                Recipient: "test@example.com",
                Subject:   "Test Subject",
                Content:   "Test Content",
            },
            wantErr: false,
        },
        {
            name: "successful telegram notification",
            notification: &model.Notification{
                Type:      model.TelegramNotification,
                Recipient: "123456789",
                Content:   "Test Content",
            },
            wantErr: false,
        },
        // Add more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := svc.SendNotification(context.Background(), tt.notification)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
} 