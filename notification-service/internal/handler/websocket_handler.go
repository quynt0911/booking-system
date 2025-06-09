package handler

import (
	"notification-service/internal/websocket"

	"github.com/gin-gonic/gin"
)

func WebSocketHandler(hub *websocket.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		websocket.ServeWs(hub, c.Writer, c.Request)
	}
}
