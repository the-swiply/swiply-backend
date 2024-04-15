package converter

import (
	"github.com/the-swiply/swiply-backend/event/internal/domain"
	"github.com/the-swiply/swiply-backend/event/pkg/api/event"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func EventFromCreateEventRequest(req *event.CreateEventRequest) domain.Event {
	return domain.Event{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Date:        req.GetDate().AsTime(),
	}
}

func EventFromUpdateEventRequest(req *event.UpdateEventRequest) domain.Event {
	return domain.Event{
		ID:          req.GetEvent().GetEventId(),
		Title:       req.GetEvent().GetTitle(),
		Description: req.GetEvent().GetDescription(),
		Date:        req.GetEvent().GetDate().AsTime(),
	}
}

func EventsToGetUserEventsResponse(events []domain.Event) *event.GetUserEventsResponse {
	res := &event.GetUserEventsResponse{Event: make([]*event.EventModel, 0, len(events))}
	for _, ev := range events {
		// TODO: add members and photos
		res.Event = append(res.Event, &event.EventModel{
			EventId:     ev.ID,
			Title:       ev.Title,
			Description: ev.Description,
			Date:        timestamppb.New(ev.Date),
		})
	}

	return res
}
