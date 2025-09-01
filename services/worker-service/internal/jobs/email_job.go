package jobs

import (
	"log"
	"worker-service/internal/service"
)

func RunEmailJob() {
	log.Println("Running Email Queue Job")
	service.ProcessQueuedEmails()
}
