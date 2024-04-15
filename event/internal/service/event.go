package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/the-swiply/swiply-backend/event/internal/domain"
	"github.com/the-swiply/swiply-backend/pkg/houston/auf"
)

type EventRepository interface {
	CreateEvent(ctx context.Context, event domain.Event) (int64, error)
	UpdateEvent(ctx context.Context, event domain.Event) error
	GetEvents(ctx context.Context, owner uuid.UUID) ([]domain.Event, error)
}

type EventService struct {
	cfg             EventConfig
	eventRepository EventRepository
}

func NewEventService(cfg EventConfig, eventRepository EventRepository) *EventService {
	return &EventService{
		cfg:             cfg,
		eventRepository: eventRepository,
	}
}

func (e *EventService) CreateEvent(ctx context.Context, event domain.Event) (int64, error) {
	userID := auf.ExtractUserIDFromContext[uuid.UUID](ctx)
	event.Owner = userID

	eventID, err := e.eventRepository.CreateEvent(ctx, event)
	if err != nil {
		return 0, err
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

func (e *EventService) GetUserEvents(ctx context.Context) ([]domain.Event, error) {
	userID := auf.ExtractUserIDFromContext[uuid.UUID](ctx)

	events, err := e.eventRepository.GetEvents(ctx, userID)
	if err != nil {
		return nil, err
	}

	return events, nil
}
