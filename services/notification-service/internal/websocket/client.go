package websocket

import "log"

type Client struct {
	ID string
}

func (c *Client) Send(message string) {
	log.Printf("Sending message to client %s: %s", c.ID, message)
}