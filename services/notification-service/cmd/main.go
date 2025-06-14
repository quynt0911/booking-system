package main

import (
	"log"
	"notification-service/internal/config"
	"notification-service/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	r := gin.Default()
	routes.SetupRoutes(r)
	log.Printf("Notification service running on port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}