package routes

import (
	"notification-service/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	notificationHandler *handler.NotificationHandler,
	websocketHandler *handler.WebSocketHandler,
	settingsHandler *handler.SettingsHandler,
) {
	// WebSocket route
	router.GET("/ws", websocketHandler.HandleWebSocket)

	// Other routes...
	router.GET("/notifications", notificationHandler.GetRecentNotifications)
	router.POST("/settings", settingsHandler.UpdateNotificationSettings)
}
