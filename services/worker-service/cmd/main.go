package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/your-org/booking-system/services/worker-service/internal/config"
	"github.com/your-org/booking-system/services/worker-service/internal/handler"
	"github.com/your-org/booking-system/services/worker-service/internal/scheduler"
	"github.com/your-org/booking-system/services/worker-service/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize scheduler
	cronScheduler := scheduler.NewCron()
	queueScheduler := scheduler.NewQueue(cfg.Redis)

	// Initialize services
	cronService := service.NewCronService(cronScheduler)
	reminderService := service.NewReminderService(cfg)
	emailQueueService := service.NewEmailQueueService(cfg)
	cleanupService := service.NewCleanupService(cfg)

	// Initialize handler
	jobHandler := handler.NewJobHandler(cronService, reminderService, emailQueueService, cleanupService)

	// Setup router
	router := gin.Default()
	router.POST("/jobs", jobHandler.CreateJob)
	router.GET("/jobs/:id", jobHandler.GetJob)
	router.DELETE("/jobs/:id", jobHandler.DeleteJob)

	// Start server
	go func() {
		if err := router.Run(":8085"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Start schedulers
	go cronScheduler.Start()
	go queueScheduler.Start()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")
	cronScheduler.Stop()
	queueScheduler.Stop()
}
