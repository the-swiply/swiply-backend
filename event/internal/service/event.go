package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/the-swiply/swiply-backend/event/internal/domain"
	"github.com/the-swiply/swiply-backend/pkg/houston/auf"
)

type EventRepository interface {
	CreateEvent(ctx context.Context, event domain.Event) (int64, error)
	UpdateEvent(ctx context.Context, event domain.Event) error
	GetEvents(ctx context.Context, limit, offset int64) ([]domain.Event, error)
	GetUserOwnEvents(ctx context.Context, owner uuid.UUID) ([]domain.Event, error)
	GetUserMembershipEvents(ctx context.Context, member uuid.UUID) ([]domain.Event, error)
	GetEventMembers(ctx context.Context, eventID int64) ([]domain.UserEventStatus, error)
	GetEventByID(ctx context.Context, id int64) (domain.Event, error)
	JoinEvent(ctx context.Context, eventID int64, userID uuid.UUID) error
	AcceptEventJoin(ctx context.Context, eventID int64, owner, userID uuid.UUID) error
}

type ChatClient interface {
	CreateChat(ctx context.Context, members []uuid.UUID) (int64, error)
	AddChatMembers(ctx context.Context, chatID int64, members []uuid.UUID) error
}

type EventService struct {
	cfg             EventConfig
	eventRepository EventRepository
	chatClient      ChatClient
}

func NewEventService(cfg EventConfig, eventRepository EventRepository, chatClient ChatClient) *EventService {
	return &EventService{
		cfg:             cfg,
		eventRepository: eventRepository,
		chatClient:      chatClient,
	}
}

func (e *EventService) CreateEvent(ctx context.Context, event domain.Event) (int64, error) {
	userID := auf.ExtractUserIDFromContext[uuid.UUID](ctx)
	event.Owner = userID

	eventID, err := e.eventRepository.CreateEvent(ctx, event)
	if err != nil {
		return 0, err
	}

	chatID, err := e.chatClient.CreateChat(ctx, []uuid.UUID{userID})
	if err != nil {
		return 0, fmt.Errorf("can't create chat: %w", err)
	}

	event.ChatID = chatID
	err = e.eventRepository.UpdateEvent(ctx, event)
	if err != nil {
		return 0, fmt.Errorf("can't set chat it for event: %w", err)
	}

	return eventID, nil
}

func (e *EventService) UpdateEvent(ctx context.Context, event domain.Event) error {
	userID := auf.ExtractUserIDFromContext[uuid.UUID](ctx)
	event.Owner = userID

	err := e.eventRepository.UpdateEvent(ctx, event)
	if err != nil {
		return err
	}

	return nil
}

func (e *EventService) GetEvents(ctx context.Context, limit, offset int64) ([]domain.Event, error) {
	events, err := e.eventRepository.GetEvents(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (e *EventService) GetUserOwnEvents(ctx context.Context) ([]domain.Event, error) {
	userID := auf.ExtractUserIDFromContext[uuid.UUID](ctx)

	events, err := e.eventRepository.GetUserOwnEvents(ctx, userID)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (e *EventService) GetUserMembershipEvents(ctx context.Context) ([]domain.Event, error) {
	userID := auf.ExtractUserIDFromContext[uuid.UUID](ctx)

	events, err := e.eventRepository.GetUserMembershipEvents(ctx, userID)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (e *EventService) GetEventMembers(ctx context.Context, eventID int64) ([]domain.UserEventStatus, error) {
	events, err := e.eventRepository.GetEventMembers(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (e *EventService) JoinEvent(ctx context.Context, eventID int64) error {
	userID := auf.ExtractUserIDFromContext[uuid.UUID](ctx)

	err := e.eventRepository.JoinEvent(ctx, eventID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (e *EventService) AcceptEventJoin(ctx context.Context, eventID int64, userIDToAdd uuid.UUID) error {
	userID := auf.ExtractUserIDFromContext[uuid.UUID](ctx)

	ev, err := e.eventRepository.GetEventByID(ctx, eventID)
	if err != nil {
		return err
	}

	err = e.eventRepository.AcceptEventJoin(ctx, eventID, userID, userIDToAdd)
	if err != nil {
		return err
	}

	err = e.chatClient.AddChatMembers(ctx, ev.ID, []uuid.UUID{userIDToAdd})
	if err != nil {
		return fmt.Errorf("can't add user to event's chat: %w", err)
	}

	return nil
}
