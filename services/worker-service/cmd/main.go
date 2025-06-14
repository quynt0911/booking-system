package main

import (
	"log"
	"worker/internal/scheduler"
)

func main() {
	log.Println("Worker Service started")
	scheduler.StartQueueProcessor()
}
