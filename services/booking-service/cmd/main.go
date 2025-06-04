package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"booking-service/internal/config"
	"booking-service/internal/handler"
	"booking-service/internal/repository"
	"booking-service/internal/routes"
	"booking-service/internal/service"
	"shared/pkg/database"
	"shared/pkg/logger"
)

func main() {
	// Initialize logger
	logger := logger.NewLogger()
	
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize database
	db, err := database.NewPostgresConnection(cfg.Database.URL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize Redis
	redisClient, err := database.NewRedisConnection(cfg.Redis.URL)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redisClient.Close()

	// Initialize repositories
	bookingRepo := repository.NewBookingRepository(db)
	statusHistoryRepo := repository.NewStatusHistoryRepository(db)

	// Initialize services
	bookingService := service.NewBookingService(bookingRepo, statusHistoryRepo, redisClient, logger)
	conflictChecker := service.NewConflictChecker(bookingRepo, redisClient)
	statusService := service.NewStatusService(statusHistoryRepo, logger)

	// Initialize handlers
	bookingHandler := handler.NewBookingHandler(bookingService, conflictChecker, logger)
	statusHandler := handler.NewStatusHandler(statusService, logger)
	historyHandler := handler.NewHistoryHandler(bookingService, logger)

	// Initialize Gin router
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	
	router := gin.Default()
	
	// Setup routes
	routes.SetupRoutes(router, bookingHandler, statusHandler, historyHandler)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.App.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info(fmt.Sprintf("Booking service starting on port %s", cfg.App.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down booking service...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	logger.Info("Booking service stopped")
}