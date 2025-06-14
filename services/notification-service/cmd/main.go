package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/your-org/booking-system/services/notification-service/internal/config"
	"github.com/your-org/booking-system/services/notification-service/internal/handler"
	"github.com/your-org/booking-system/services/notification-service/internal/routes"
	"github.com/your-org/booking-system/services/notification-service/internal/service"
	"github.com/your-org/booking-system/services/notification-service/internal/websocket"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize websocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize services
	notificationService := service.NewNotificationService(cfg)
	emailService := service.NewEmailService(cfg)
	telegramService := service.NewTelegramService(cfg)
	websocketService := service.NewWebSocketService(hub)

	// Initialize handlers
	notificationHandler := handler.NewNotificationHandler(notificationService)
	websocketHandler := handler.NewWebSocketHandler(websocketService)
	settingsHandler := handler.NewSettingsHandler(notificationService)

	// Setup router
	router := gin.Default()
	routes.SetupRoutes(router, notificationHandler, websocketHandler, settingsHandler)

	// Start server
	go func() {
		if err := router.Run(":8084"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
