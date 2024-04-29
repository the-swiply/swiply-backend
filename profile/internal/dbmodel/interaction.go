package dbmodel

import (
	"time"

	"github.com/google/uuid"
)

type Interaction struct {
	ID        int64
	From      uuid.UUID
	To        uuid.UUID
	Type      string
	CreatedAt time.Time
}
