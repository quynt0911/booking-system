package scheduler

import (
	"log"
	"worker/internal/jobs"
)

func StartQueueProcessor() {
	log.Println("Starting worker queue processor...")

	jobs.RunReminderJob()
	jobs.RunEmailJob()
	jobs.RunCleanupJob()
	jobs.RunStatusUpdateJob()
}