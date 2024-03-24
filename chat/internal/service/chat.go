package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/the-swiply/swiply-backend/chat/internal/domain"
	"github.com/the-swiply/swiply-backend/pkg/houston/auf"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"
	"slices"
	"sync"
	"time"
)

type SequenceGenerator interface {
	GenerateNextID(ctx context.Context, chatID int64) (int64, error)
	RollbackID(ctx context.Context, chatID int64) error
}

type ChatRepository interface {
	GetChatMembers(ctx context.Context, chatID int64) ([]uuid.UUID, error)
	SaveMessage(ctx context.Context, msg domain.ChatMessage) error
	GetNextMessages(ctx context.Context, chatID int64, start int64, limit int64) ([]domain.ChatMessage, error)
	GetPreviousMessages(ctx context.Context, chatID int64, start int64, limit int64) ([]domain.ChatMessage, error)
	GetUserChats(ctx context.Context, userID uuid.UUID) ([]domain.Chat, error)
	RemoveUserFromChatMembers(ctx context.Context, userID uuid.UUID, chatID int64) error
	CreateChat(ctx context.Context, members []uuid.UUID) (int64, error)
}

type MessagePublisher interface {
	PublishMessage(ctx context.Context, msg domain.Message) error
}

type ChatClient interface {
	SendMessage(msg domain.Message) error
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
	userID := auf.ExtractUserIDFromContext[uuid.UUID](ctx)

	err := c.checkUserInChatWithError(ctx, userID, chatID)
	if err != nil {
		return err
	}

	gmu := c.createChatLock(chatID)
	err = gmu.Lock(ctx)
	defer func() {
		err := gmu.Unlock(ctx)
		if err != nil {
			loggy.Errorf("unlock chat error: %v", err)
		}
	}()

	idInChat, err := c.seqGen.GenerateNextID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("can't generate message sequence id for chat: %w", err)
	}

	chatMsg := domain.ChatMessage{
		ID:       uuid.New(),
		From:     userID,
		ChatID:   chatID,
		IDInChat: idInChat,
		SendTime: time.Now(),
		Content:  content,
	}

	err = c.chatRepository.SaveMessage(ctx, chatMsg)
	if err != nil {
		rollbackErr := c.seqGen.RollbackID(ctx, chatID)
		if rollbackErr != nil {
			loggy.Errorf("can't rollback message sequence id for chat: %v", err)
		}

		return fmt.Errorf("can't save message: %w", err)
	}

	msg := domain.Message{
		Type:    domain.MessageTypeMessage,
		ChatID:  chatID,
		Payload: chatMsg,
	}

	err = c.messagePublisher.PublishMessage(ctx, msg)
	if err != nil {
		loggy.Errorf("can't publish message on receive chat message: %v", err)
	}

	return nil
}

func (c *ChatService) SendMessageToChat(ctx context.Context, msg domain.Message) error {
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

func (c *ChatService) GetNextMessages(ctx context.Context, chatID int64, start int64, limit int64) ([]domain.ChatMessage, error) {
	userID := auf.ExtractUserIDFromContext[uuid.UUID](ctx)

	err := c.checkUserInChatWithError(ctx, userID, chatID)
	if err != nil {
		return nil, err
	}

	messages, err := c.chatRepository.GetNextMessages(ctx, chatID, start, limit)
	if err != nil {
		return nil, fmt.Errorf("can't get messages: %w", err)
	}

	return messages, nil
}

func (c *ChatService) GetPreviousMessages(ctx context.Context, chatID int64, start int64, limit int64) ([]domain.ChatMessage, error) {
	userID := auf.ExtractUserIDFromContext[uuid.UUID](ctx)

	err := c.checkUserInChatWithError(ctx, userID, chatID)
	if err != nil {
		return nil, err
	}

	messages, err := c.chatRepository.GetPreviousMessages(ctx, chatID, start, limit)
	if err != nil {
		return nil, fmt.Errorf("can't get messages: %w", err)
	}

	return messages, nil
}

func (c *ChatService) GetUserChats(ctx context.Context) ([]domain.Chat, error) {
	userID := auf.ExtractUserIDFromContext[uuid.UUID](ctx)

	return c.chatRepository.GetUserChats(ctx, userID)
}

func (c *ChatService) LeaveChat(ctx context.Context, chatID int64) error {
	userID := auf.ExtractUserIDFromContext[uuid.UUID](ctx)
	err := c.checkUserInChatWithError(ctx, userID, chatID)
	if err != nil {
		return err
	}

	err = c.chatRepository.RemoveUserFromChatMembers(ctx, userID, chatID)
	if err != nil {
		return err
	}

	msg := domain.Message{
		Type:   domain.MessageTypeUserLeft,
		ChatID: chatID,
		Payload: map[string]uuid.UUID{
			"user_id": userID,
		},
	}

	err = c.messagePublisher.PublishMessage(ctx, msg)
	if err != nil {
		loggy.Errorf("can't publish message on leave chat: %v", err)
	}

	return nil
}

func (c *ChatService) CreateChat(ctx context.Context, members []uuid.UUID) error {
	chatID, err := c.chatRepository.CreateChat(ctx, members)
	if err != nil {
		return err
	}

	msg := domain.Message{
		Type:   domain.MessageTypeChatCreated,
		ChatID: chatID,
		Payload: map[string][]uuid.UUID{
			"members": members,
		},
	}

	err = c.messagePublisher.PublishMessage(ctx, msg)
	if err != nil {
		loggy.Errorf("can't publish message on create chat: %v", err)
	}

	return nil
}

func (c *ChatService) checkUserInChatWithError(ctx context.Context, userID uuid.UUID, chatID int64) error {
	chatMembers, err := c.chatRepository.GetChatMembers(ctx, chatID)

	if errors.Is(err, domain.ErrEntityIsNotExists) {
		return domain.ErrEntityIsNotExists
	}
	if err != nil {
		return fmt.Errorf("can't get chat members: %w", err)
	}

	if !slices.Contains(chatMembers, userID) {
		return domain.ErrUserNotInChat
	}

	return nil
}
