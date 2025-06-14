package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/your-org/booking-system/services/notification-service/internal/service"
)

type WebSocketHandler struct {
	websocketService *service.WebSocketService
}

func NewWebSocketHandler(websocketService *service.WebSocketService) *WebSocketHandler {
	return &WebSocketHandler{
		websocketService: websocketService,
	}
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	h.websocketService.manager.HandleWebSocket(c)
}
