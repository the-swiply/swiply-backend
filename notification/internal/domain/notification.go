package domain

import (
	"github.com/google/uuid"
)

type Notification struct {
	ID          uuid.UUID
	DeviceToken string
}
