package dbmodel

import (
	"time"

	"github.com/google/uuid"
)

type Location struct {
	Lat  float64
	Long float64
}

type Profile struct {
	ID           uuid.UUID
	Email        string
	Name         string
	Interests    []int64
	BirthDay     time.Time
	Gender       string
	Lat          float64 `db:"location_lat"`
	Long         float64 `db:"location_long"`
	Info         string
	Subscription string
	UpdatedAt    time.Time
}
