package jobs

import (
	"log"
	"worker/internal/service"
)

func RunCleanupJob() {
	log.Println("Running Cleanup Job")
	service.CleanupOldRecords()
}