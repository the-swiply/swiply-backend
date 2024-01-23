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

type ChatLocker interface {
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) error
}

type ChatService struct {
	cfg              ChatConfig
	seqGen           SequenceGenerator
	chatRepository   ChatRepository
	messagePublisher MessagePublisher
	createChatLock   func(chatID int64) ChatLocker

	chatClients   map[uuid.UUID]ChatClient
	chatClientsMu sync.RWMutex
}

func NewChatService(cfg ChatConfig, seqGen SequenceGenerator, chatRepository ChatRepository,
	messagePublisher MessagePublisher, createChatLock func(chatID int64) ChatLocker) *ChatService {
	return &ChatService{
		cfg:              cfg,
		seqGen:           seqGen,
		chatRepository:   chatRepository,
		messagePublisher: messagePublisher,
		createChatLock:   createChatLock,
		chatClients:      make(map[uuid.UUID]ChatClient),
		chatClientsMu:    sync.RWMutex{},
	}
}

func (c *ChatService) ReceiveChatMessage(ctx context.Context, chatID int64, content string) error {
	userID := extractUserIDFromContext(ctx)

	ok, err := c.checkUserInChat(ctx, userID, chatID)
	if err != nil {
		return fmt.Errorf("can't check if user in chat: %w", err)
	}
	if !ok {
		return domain.ErrUserNotInChat
	}

	mu := c.createChatLock(chatID)
	err = mu.Lock(ctx)
	defer func() {
		err := mu.Unlock(ctx)
		if err != nil {
			loggy.Errorf("unlock chat error: %v", err)
		}
	}()

	idInChat, err := c.seqGen.GenerateNextID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("can't generate sequence id for chat: %w", err)
	}

	msg := domain.ChatMessage{
		ID:       uuid.New(),
		From:     userID,
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
		loggy.Errorf("can't get chat members: %v", err)
		return nil
	}

	c.chatClientsMu.RLock()
	defer c.chatClientsMu.RUnlock()

	for _, member := range chatMembers {
		if client, ok := c.chatClients[member]; ok {
			err = client.SendMessage(msg)
			if err != nil {
				loggy.Errorln(err)
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

func (c *ChatService) checkUserInChat(ctx context.Context, userID uuid.UUID, chatID int64) (bool, error) {
	chatMembers, err := c.chatRepository.GetChatMembers(ctx, chatID)
	if err != nil {
		return false, fmt.Errorf("can't get chat members: %w", err)
	}

	return slices.Contains(chatMembers, userID), nil
}
