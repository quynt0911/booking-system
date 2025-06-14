package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/your-org/booking-system/services/notification-service/internal/handler"
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
