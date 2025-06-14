package websocket

type Hub struct {
	Clients map[string]*Client
}

func NewHub() *Hub {
	return &Hub{Clients: make(map[string]*Client)}
}

func (h *Hub) Broadcast(userID, message string) {
	if c, ok := h.Clients[userID]; ok {
		c.Send(message)
	}
}