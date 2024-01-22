package sevents

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/the-swiply/swiply-backend/chat/internal/domain"
)

type Client struct {
	conn *websocket.Conn
}

func (c *Client) SendMessage(msg domain.ChatMessage) error {
	err := c.conn.WriteJSON(msg)
	if errors.Is(err, websocket.ErrCloseSent) {
		return nil
	}

	return err
}
