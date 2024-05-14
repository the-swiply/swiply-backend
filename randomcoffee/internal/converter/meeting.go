package converter

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/the-swiply/swiply-backend/randomcoffee/internal/domain"
	"github.com/the-swiply/swiply-backend/randomcoffee/pkg/api/randomcoffee"
)

func MeetingFromDomainToProto(meeting domain.Meeting) *randomcoffee.Meeting {
	meet := &randomcoffee.Meeting{
		Id:             meeting.ID.String(),
		OwnerId:        meeting.OwnerID.String(),
		MemberId:       meeting.MemberID.String(),
		Start:          timestamppb.New(meeting.Start),
		End:            timestamppb.New(meeting.End),
		OrganizationId: meeting.OrganizationID,
		CreatedAt:      timestamppb.New(meeting.CreatedAt),
	}

	switch meeting.Status {
	case domain.MeetingStatusUnspecified:
		meet.Status = randomcoffee.MeetingStatus_MEETING_STATUS_UNSPECIFIED
	case domain.MeetingStatusAwaitingSchedule:
		meet.Status = randomcoffee.MeetingStatus_AWAITING_SCHEDULE
	case domain.MeetingStatusScheduling:
		meet.Status = randomcoffee.MeetingStatus_SCHEDULING
	case domain.MeetingStatusScheduled:
		meet.Status = randomcoffee.MeetingStatus_SCHEDULED
	}

	return meet
}

func MeetingFromProtoToDomain(meeting *randomcoffee.Meeting) (domain.Meeting, error) {
	var (
		id       uuid.UUID
		ownerID  uuid.UUID
		memberID uuid.UUID
		err      error
	)

	id, err = uuid.Parse(meeting.Id)
	if err != nil {
		return domain.Meeting{}, fmt.Errorf("can't parse meeting id: %w", err)
	}

	ownerID, err = uuid.Parse(meeting.OwnerId)
	if err != nil {
		return domain.Meeting{}, fmt.Errorf("can't parse owner id: %w", err)
	}

	if meeting.MemberId != "" {
		memberID, err = uuid.Parse(meeting.MemberId)
		if err != nil {
			return domain.Meeting{}, fmt.Errorf("can't parse member id: %w", err)
		}
	}

	meet := domain.Meeting{
		ID:             id,
		OwnerID:        ownerID,
		MemberID:       memberID,
		Start:          meeting.Start.AsTime(),
		End:            meeting.End.AsTime(),
		OrganizationID: meeting.OrganizationId,
		CreatedAt:      meeting.CreatedAt.AsTime(),
	}

	switch meeting.Status {
	case randomcoffee.MeetingStatus_MEETING_STATUS_UNSPECIFIED:
		meet.Status = domain.MeetingStatusUnspecified
	case randomcoffee.MeetingStatus_AWAITING_SCHEDULE:
		meet.Status = domain.MeetingStatusAwaitingSchedule
	case randomcoffee.MeetingStatus_SCHEDULING:
		meet.Status = domain.MeetingStatusScheduling
	case randomcoffee.MeetingStatus_SCHEDULED:
		meet.Status = domain.MeetingStatusScheduled
	}

	return meet, nil
}
