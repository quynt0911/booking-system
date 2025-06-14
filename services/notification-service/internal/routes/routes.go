package routes

import (
	"github.com/gin-gonic/gin"
	"notification-service/internal/handler"
	"notification-service/internal/websocket"
)

func SetupRoutes(r *gin.Engine) {
	hub := websocket.NewHub()
	r.GET("/ws", handler.WebSocketHandler(hub))
	r.GET("/notifications", handler.GetRecentNotifications)
	r.POST("/settings", handler.UpdateNotificationSettings)
}