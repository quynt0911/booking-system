package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/booking-system/services/notification-service/internal/websocket"
)

func TestNotificationFlow(t *testing.T) {
    // Setup WebSocket connection
    ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:8084/ws", nil)
    assert.NoError(t, err)
    defer ws.Close()

    // Test cases
    tests := []struct {
        name     string
        message  string
        wantResp string
    }{
        {
            name:     "receive booking notification",
            message:  `{"type":"booking","content":"New booking request"}`,
            wantResp: `{"type":"booking","status":"received"}`,
        },
        // Add more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Send message
            err := ws.WriteMessage(websocket.TextMessage, []byte(tt.message))
            assert.NoError(t, err)

            // Read response
            _, message, err := ws.ReadMessage()
            assert.NoError(t, err)
            assert.Equal(t, tt.wantResp, string(message))
        })
    }
} 