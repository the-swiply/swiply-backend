package domain

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	ID          int64
	Owner       uuid.UUID
	Members     []uuid.UUID
	Title       string
	Description string
	Date        time.Time
}
