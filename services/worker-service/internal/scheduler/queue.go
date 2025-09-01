package scheduler

import (
	"log"
	"worker-service/internal/jobs"
)

func StartQueueProcessor() {
	log.Println("Starting worker queue processor...")

	jobs.RunReminderJob()
	jobs.RunEmailJob()
	jobs.RunCleanupJob()
	jobs.RunStatusUpdateJob()
}
