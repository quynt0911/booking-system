package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/your-org/booking-system/services/notification-service/internal/handler"
)

func SetupRoutes(
	router *gin.Engine,
	notificationHandler *handler.NotificationHandler,
	websocketHandler interface{},
	settingsHandler interface{},
) {
	if websocketHandler != nil {
		// router.GET("/ws", websocketHandler.HandleWebSocket)
	}
	if notificationHandler != nil {
		router.POST("/send-email", notificationHandler.SendEmailHandler)
	}
	// Các route khác nếu cần
}
