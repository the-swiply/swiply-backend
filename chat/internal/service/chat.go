package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/the-swiply/swiply-backend/chat/internal/domain"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"
	"slices"
	"sync"
	"time"
)

type SequenceGenerator interface {
	GenerateNextID(ctx context.Context, chatID int64) (int64, error)
}

type ChatRepository interface {
	GetChatMembers(ctx context.Context, chatID int64) ([]uuid.UUID, error)
	SaveMessage(ctx context.Context, msg domain.ChatMessage) error
}

type MessagePublisher interface {
	PublishMessage(ctx context.Context, msg domain.ChatMessage) error
}

type ChatClient interface {
	SendMessage(msg domain.ChatMessage) error
}

type ChatService struct {
	cfg              ChatConfig
	seqGen           SequenceGenerator
	chatRepository   ChatRepository
	messagePublisher MessagePublisher

	chatClients   map[uuid.UUID]ChatClient
	chatClientsMu sync.RWMutex
}

func NewChatService(cfg ChatConfig, seqGen SequenceGenerator, chatRepository ChatRepository,
	messagePublisher MessagePublisher) *ChatService {
	return &ChatService{
		cfg:              cfg,
		seqGen:           seqGen,
		chatRepository:   chatRepository,
		messagePublisher: messagePublisher,
		chatClients:      make(map[uuid.UUID]ChatClient),
		chatClientsMu:    sync.RWMutex{},
	}
}

func (c *ChatService) ReceiveChatMessage(ctx context.Context, chatID int64, content string) error {
	ok, err := c.checkUserInChat(ctx, chatID)
	if err != nil {
		return fmt.Errorf("can't check if user in chat: %w", err)
	}
	if !ok {
		return domain.ErrUserNotInChat
	}

	idInChat, err := c.seqGen.GenerateNextID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("can't generate sequence id for chat: %w", err)
	}

	msg := domain.ChatMessage{
		ID:       uuid.New(),
		From:     uuid.New(),
		ChatID:   chatID,
		IDInChat: idInChat,
		SendTime: time.Now(),
		Content:  content,
	}

	err = c.chatRepository.SaveMessage(ctx, msg)
	if err != nil {
		return fmt.Errorf("can't save message: %w", err)
	}

	err = c.messagePublisher.PublishMessage(ctx, msg)
	if err != nil {
		loggy.Errorf("can't publish message: %v", err)
	}

	return nil
}

func (c *ChatService) SendChatMessage(ctx context.Context, msg domain.ChatMessage) error {
	chatMembers, err := c.chatRepository.GetChatMembers(ctx, msg.ChatID)
	if err != nil && !errors.Is(err, domain.ErrEntityIsNotExists) {
		err = fmt.Errorf("can't get chat members: %w", err)
		loggy.Errorln(err)
		return err
	}

	c.chatClientsMu.RLock()
	defer c.chatClientsMu.RUnlock()

	for _, member := range chatMembers {
		if client, ok := c.chatClients[member]; ok {
			err = client.SendMessage(msg)
			if err != nil {
				loggy.Infoln(err)
			}
		}
	}

	return nil
}

func (c *ChatService) AddChatClient(userID uuid.UUID, client ChatClient) {
	c.chatClientsMu.Lock()
	c.chatClients[userID] = client
	c.chatClientsMu.Unlock()
}

func (c *ChatService) RemoveChatClient(userID uuid.UUID) {
	c.chatClientsMu.Lock()
	delete(c.chatClients, userID)
	c.chatClientsMu.Unlock()
}

func (c *ChatService) checkUserInChat(ctx context.Context, chatID int64) (bool, error) {
	userID := extractUserIDFromContext(ctx)

	chatMembers, err := c.chatRepository.GetChatMembers(ctx, chatID)
	if err != nil {
		return false, fmt.Errorf("can't get chat members: %w", err)
	}

	return slices.Contains(chatMembers, userID), nil
}
