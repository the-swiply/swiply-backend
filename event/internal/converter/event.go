package converter

import (
	"github.com/google/uuid"
	"github.com/the-swiply/swiply-backend/event/internal/domain"
	"github.com/the-swiply/swiply-backend/event/pkg/api/event"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func EventFromCreateEventRequest(req *event.CreateEventRequest) domain.Event {
	photos := make([]domain.Photo, 0, len(req.GetPhotos()))
	for _, phContent := range req.GetPhotos() {
		photos = append(photos, domain.Photo{
			Content: phContent,
		})
	}

	return domain.Event{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Photos:      photos,
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

func EventsToGetEventsResponse(events []domain.Event) *event.GetEventsResponse {
	res := &event.GetEventsResponse{Events: make([]*event.EventModel, 0, len(events))}
	for _, ev := range events {
		photoContents := make([][]byte, 0, len(ev.Photos))
		for _, ph := range ev.Photos {
			photoContents = append(photoContents, ph.Content)
		}

		res.Events = append(res.Events, &event.EventModel{
			EventId:     ev.ID,
			Title:       ev.Title,
			Description: ev.Description,
			ChatId:      ev.ChatID,
			Photos:      photoContents,
			Date:        timestamppb.New(ev.Date),
		})
	}

	return res
}

func EventsToGetUserMembershipEventsResponse(events []domain.Event) *event.GetUserMembershipEventsResponse {
	res := &event.GetUserMembershipEventsResponse{Events: make([]*event.EventModel, 0, len(events))}
	for _, ev := range events {
		photoContents := make([][]byte, 0, len(ev.Photos))
		for _, ph := range ev.Photos {
			photoContents = append(photoContents, ph.Content)
		}

		res.Events = append(res.Events, &event.EventModel{
			EventId:     ev.ID,
			Title:       ev.Title,
			Description: ev.Description,
			Photos:      photoContents,
			ChatId:      ev.ChatID,
			Date:        timestamppb.New(ev.Date),
		})
	}

	return res
}

func EventsToGetUserOwnEvents(events []domain.Event) *event.GetUserOwnEventsResponse {
	res := &event.GetUserOwnEventsResponse{Events: make([]*event.EventModel, 0, len(events))}
	for _, ev := range events {
		photoContents := make([][]byte, 0, len(ev.Photos))
		for _, ph := range ev.Photos {
			photoContents = append(photoContents, ph.Content)
		}

		res.Events = append(res.Events, &event.EventModel{
			EventId:     ev.ID,
			Title:       ev.Title,
			Description: ev.Description,
			Photos:      photoContents,
			ChatId:      ev.ChatID,
			Date:        timestamppb.New(ev.Date),
		})
	}

	return res
}

func UserEventStatusToPB(members []domain.UserEventStatus) *event.GetEventMembersResponse {
	res := &event.GetEventMembersResponse{UsersStatuses: make([]*event.GetEventMembersResponse_UserWithEventStatus, 0, len(members))}

	for _, member := range members {
		var st event.UserEventStatus
		switch member.Status {
		case "join_request":
			st = event.UserEventStatus_JOIN_REQUEST
		case "member":
			st = event.UserEventStatus_MEMBER
		default:
			st = event.UserEventStatus_USER_EVENT_STATUS_UNKNOWN
		}

		res.UsersStatuses = append(res.UsersStatuses, &event.GetEventMembersResponse_UserWithEventStatus{
			UserId: member.UserID.String(),
			Status: st,
		})
	}

	return res
}

func UUIDsToStrings(uuids []uuid.UUID) []string {
	res := make([]string, 0, len(uuids))

	for _, u := range uuids {
		res = append(res, u.String())
	}

	return res
}
