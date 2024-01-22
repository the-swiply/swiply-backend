package domain

import (
	"github.com/google/uuid"
	"time"
)

type ChatMessage struct {
	ID       uuid.UUID `json:"id"`
	From     uuid.UUID `json:"from"`
	ChatID   int64     `json:"chat_id"`
	IDInChat int64     `json:"id_in_chat"`
	SendTime time.Time `json:"send_time"`
	Content  string    `json:"content"`
}
