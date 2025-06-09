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

	"booking-system/services/booking-service/internal/config"
	"booking-system/services/booking-service/internal/handler"
	"booking-system/services/booking-service/internal/repository"
	"booking-system/services/booking-service/internal/routes"
	"booking-system/services/booking-service/internal/service"
	"booking-system/services/booking-service/pkg/database"
	"booking-system/services/booking-service/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logger
	appLogger := logger.NewLogger()

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

	// Convert *sql.DB to *gorm.DB
	gormDB, err := database.GetGormDB(db)
	if err != nil {
		log.Fatal("Failed to get GORM DB:", err)
	}

	// Initialize Redis
	redisClient, err := database.NewRedisConnection(cfg.Redis.URL)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redisClient.Close()

	// Initialize repositories
	bookingRepo := repository.NewBookingRepository(gormDB)
	statusHistoryRepo := repository.NewStatusHistoryRepository(gormDB)

	// Initialize services
	bookingService := service.NewBookingService(bookingRepo, statusHistoryRepo, redisClient, appLogger)
	conflictChecker := service.NewConflictChecker(bookingRepo, redisClient)
	statusService := service.NewStatusService(statusHistoryRepo, bookingRepo, appLogger)

	// Initialize handlers
	bookingHandler := handler.NewBookingHandler(bookingService, conflictChecker, appLogger)
	statusHandler := handler.NewStatusHandler(statusService, appLogger)
	historyHandler := handler.NewHistoryHandler(bookingService, appLogger)

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
		appLogger.Info(fmt.Sprintf("Booking service starting on port %s", cfg.App.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down booking service...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	appLogger.Info("Booking service stopped")
}
