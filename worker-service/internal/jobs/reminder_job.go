package jobs

import (
	"log"
	"worker/internal/service"
)

func RunReminderJob() {
	log.Println("Running Reminder Job")
	service.SendReminder("123", "15:00")
}