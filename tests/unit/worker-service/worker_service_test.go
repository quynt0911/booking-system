package worker_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/booking-system/services/worker-service/internal/model"
	"github.com/your-org/booking-system/services/worker-service/internal/service"
)

func TestWorkerService_ProcessJob(t *testing.T) {
    // Setup
    mockRepo := newMockJobRepository()
    mockReminderService := newMockReminderService()
    mockEmailQueueService := newMockEmailQueueService()
    mockCleanupService := newMockCleanupService()

    svc := service.NewWorkerService(mockRepo, mockReminderService, mockEmailQueueService, mockCleanupService)

    tests := []struct {
        name    string
        job     *model.Job
        wantErr bool
    }{
        {
            name: "successful reminder job",
            job: &model.Job{
                Type:     model.ReminderJob,
                Payload:  `{"booking_id": 1, "reminder_time": "2024-03-20T10:00:00Z"}`,
                Schedule: "*/5 * * * *",
            },
            wantErr: false,
        },
        // Add more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := svc.ProcessJob(context.Background(), tt.job)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
} 