package service

import (
	"context"
)

type ChatService struct {
}

func NewChatService() *ChatService {
	return &ChatService{}
}

func (c *ChatService) SendMessageToUser(ctx context.Context) error {
	return nil
}
