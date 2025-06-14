package jobs

import (
	"log"
	"worker-service/internal/service"
)

func RunCleanupJob() {
	log.Println("Running Cleanup Job")
	service.CleanupOldRecords()
}
