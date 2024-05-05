package domain

import (
	"github.com/google/uuid"
	"time"
)

type Profile struct {
	ID               uuid.UUID
	Interests        []int64
	BirthDay         time.Time
	Gender           string
	Info             string
	SubscriptionType string
	LocationLat      float64
	LocationLon      float64
}

type Interaction struct {
	From     uuid.UUID
	To       uuid.UUID
	Positive bool
}
