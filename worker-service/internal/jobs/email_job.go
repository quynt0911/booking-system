package jobs

import (
	"log"
	"worker/internal/service"
)

func RunEmailJob() {
	log.Println("Running Email Queue Job")
	service.ProcessQueuedEmails()
}