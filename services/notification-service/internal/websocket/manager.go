package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Trong môi trường production, cần kiểm tra origin kỹ hơn
	},
}

type Manager struct {
	hub *Hub
}

func NewManager(hub *Hub) *Manager {
	return &Manager{
		hub: hub,
	}
}

func (m *Manager) HandleWebSocket(c *gin.Context) {
	userID := c.GetString("user_id") // Lấy từ middleware authentication

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &Client{
		hub:    m.hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
	}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (m *Manager) SendToUser(userID string, message []byte) {
	m.hub.mu.RLock()
	defer m.hub.mu.RUnlock()

	for client := range m.hub.clients {
		if client.userID == userID {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(m.hub.clients, client)
			}
		}
	}
}

func (m *Manager) Broadcast(message []byte) {
	m.hub.Broadcast(message)
}
