package domain

import (
	"github.com/google/uuid"
	"time"
)

type Recommendation uuid.UUID

type Profile struct {
	ID        uuid.UUID
	UpdatedAt time.Time
}

type Interaction struct {
	ID        uuid.UUID
	From      uuid.UUID
	To        uuid.UUID
	Positive  bool
	UpdatedAt time.Time
}
