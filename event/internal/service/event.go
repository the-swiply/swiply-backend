package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/the-swiply/swiply-backend/event/internal/domain"
	"github.com/the-swiply/swiply-backend/pkg/houston/auf"
	"github.com/the-swiply/swiply-backend/pkg/houston/dobby"
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

type EventPhotoManager interface {
	ListPhotos(ctx context.Context, eventID int64) ([]domain.Photo, error)
	GetPhoto(ctx context.Context, eventID int64, photoID int64) (domain.Photo, error)
	UploadPhoto(ctx context.Context, eventID int64, photo domain.Photo) error
}

type ChatClient interface {
	CreateChat(ctx context.Context, members []uuid.UUID) (int64, error)
	AddChatMembers(ctx context.Context, chatID int64, members []uuid.UUID) error
}

type Transactor interface {
	WithinTransaction(ctx context.Context, f func(ctx context.Context) error, opts dobby.TxOptions) error
}

type EventService struct {
	cfg             EventConfig
	eventRepository EventRepository
	transactor      Transactor
	photoMgr        EventPhotoManager
	chatClient      ChatClient
}

func NewEventService(cfg EventConfig, eventRepository EventRepository, transactor Transactor,
	photoMgr EventPhotoManager, chatClient ChatClient) *EventService {
	return &EventService{
		cfg:             cfg,
		eventRepository: eventRepository,
		transactor:      transactor,
		photoMgr:        photoMgr,
		chatClient:      chatClient,
	}
}

func (e *EventService) CreateEvent(ctx context.Context, event domain.Event) (int64, error) {
	userID := auf.ExtractUserIDFromContext[uuid.UUID](ctx)
	event.Owner = userID

	err := e.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		eventID, err := e.eventRepository.CreateEvent(txCtx, event)
		if err != nil {
			return err
		}

		for i, photo := range event.Photos {
			photo.ID = int64(i)
			err = e.photoMgr.UploadPhoto(txCtx, eventID, photo)
			if err != nil {
				return fmt.Errorf("can't upload photo: %w", err)
			}
		}

		chatID, err := e.chatClient.CreateChat(txCtx, []uuid.UUID{userID})
		if err != nil {
			return fmt.Errorf("can't create chat: %w", err)
		}

		event.ChatID = chatID
		err = e.eventRepository.UpdateEvent(txCtx, event)
		if err != nil {
			return fmt.Errorf("can't set chat it for event: %w", err)
		}
		event.ID = eventID

		return nil
	}, dobby.TxOptions{})
	if err != nil {
		return 0, err
	}

	return event.ID, nil
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

	err = e.getPhotosForEvents(ctx, events)
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

	err = e.getPhotosForEvents(ctx, events)
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

	err = e.getPhotosForEvents(ctx, events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (e *EventService) getPhotosForEvents(ctx context.Context, events []domain.Event) error {
	for i, ev := range events {
		photos, err := e.photoMgr.ListPhotos(ctx, ev.ID)
		if err != nil {
			return fmt.Errorf("can't get photos for event %d: %w", ev.ID, err)
		}
		events[i].Photos = photos
	}

	return nil
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

	err := e.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		ev, err := e.eventRepository.GetEventByID(txCtx, eventID)
		if err != nil {
			return err
		}

		err = e.eventRepository.AcceptEventJoin(txCtx, eventID, userID, userIDToAdd)
		if err != nil {
			return err
		}

		err = e.chatClient.AddChatMembers(txCtx, ev.ID, []uuid.UUID{userIDToAdd})
		if err != nil {
			return fmt.Errorf("can't add user to event's chat: %w", err)
		}

		return nil
	}, dobby.TxOptions{})

	return err
}
