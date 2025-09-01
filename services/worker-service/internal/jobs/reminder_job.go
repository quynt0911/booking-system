package jobs

import (
	"context"
	"encoding/json"
	"log"
	"worker-service/internal/service"

	"github.com/your-org/booking-system/services/worker-service/internal/model"
)

type ReminderJob struct {
	reminderService *service.ReminderService
}

func NewReminderJob(reminderService *service.ReminderService) *ReminderJob {
	return &ReminderJob{
		reminderService: reminderService,
	}
}

func (j *ReminderJob) Execute(ctx context.Context, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	var reminderPayload model.Reminder
	if err := json.Unmarshal(data, &reminderPayload); err != nil {
		return err
	}

	return j.reminderService.SendReminder(ctx, &reminderPayload)
}

func (j *ReminderJob) GetName() string {
	return string(model.ReminderJob)
}

func RunReminderJob() {
	log.Println("Running Reminder Job")
	service.SendReminder("123", "15:00")
}
