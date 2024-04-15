package service

type EventRepository interface {
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
