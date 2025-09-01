package service

import (
	"encoding/json"

	"github.com/your-org/booking-system/services/notification-service/internal/websocket"
)

type WebSocketService struct {
	manager *websocket.Manager
}

func NewWebSocketService(hub *websocket.Hub) *WebSocketService {
	manager := websocket.NewManager(hub)
	return &WebSocketService{
		manager: manager,
	}
}

func (s *WebSocketService) Send(userID string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	s.manager.SendToUser(userID, data)
	return nil
}

func (s *WebSocketService) Broadcast(message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	s.manager.Broadcast(data)
	return nil
}

func SendRealTimeNotification(userID, message string) {
	// Publish to WebSocket or Ably (abstracted)
}
