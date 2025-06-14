package main

import (
	"log"
	"worker-service/internal/scheduler"
)

func main() {
	log.Println("Worker Service started")
	scheduler.StartQueueProcessor()
}
