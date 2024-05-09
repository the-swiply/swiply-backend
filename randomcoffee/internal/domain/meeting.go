package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type MeetingStatus = meetingStatus

type meetingStatus string

func (m *meetingStatus) Set(value string) error {
	switch value {
	case string(MeetingStatusUnspecified):
		*m = MeetingStatusUnspecified
	case string(MeetingStatusAwaitingSchedule):
		*m = MeetingStatusAwaitingSchedule
	case string(MeetingStatusScheduling):
		*m = MeetingStatusScheduling
	case string(MeetingStatusScheduled):
		*m = MeetingStatusScheduled
	default:
		return errors.New("unknown meeting status")
	}

	return nil
}

const (
	MeetingStatusUnspecified      meetingStatus = ""
	MeetingStatusAwaitingSchedule meetingStatus = "AWAITING_SCHEDULE"
	MeetingStatusScheduling       meetingStatus = "SCHEDULING"
	MeetingStatusScheduled        meetingStatus = "SCHEDULED"
)

type Meeting struct {
	ID             uuid.UUID
	OwnerID        uuid.UUID
	MemberID       uuid.UUID
	Start          time.Time
	End            time.Time
	OrganizationID int64
	Status         meetingStatus
}
