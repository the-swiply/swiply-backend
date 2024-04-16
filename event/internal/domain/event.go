package domain

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	ID          int64
	Owner       uuid.UUID
	Title       string
	Description string
	ChatID      int64
	Date        time.Time
}

type UserEventStatus struct {
	UserID  uuid.UUID
	EventID int64
	Status  string
}
