package domain

import (
	"github.com/google/uuid"
	"time"
)

const (
	MessageTypeChatCreated MessageType = "CHAT_CREATED"
	MessageTypeMessage     MessageType = "MESSAGE"
	MessageTypeUserLeft    MessageType = "USER_LEFT"
)

type MessageType string

type Message struct {
	Type    MessageType `json:"type"`
	ChatID  int64       `json:"chat_id"`
	Payload any         `json:"payload"`
}

type ChatMessage struct {
	ID       uuid.UUID `json:"id"`
	From     uuid.UUID `json:"from"`
	ChatID   int64     `json:"chat_id"`
	IDInChat int64     `json:"id_in_chat"`
	SendTime time.Time `json:"send_time"`
	Content  string    `json:"content"`
}
