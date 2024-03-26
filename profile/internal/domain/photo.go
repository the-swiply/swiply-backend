package domain

import (
	"github.com/google/uuid"
)

type Photo struct {
	ID      uuid.UUID
	Content []byte
}
