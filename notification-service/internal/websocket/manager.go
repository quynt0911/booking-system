package websocket

import "net/http"

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP to WS and register client
}