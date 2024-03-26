package dbmodel

import (
	"github.com/google/uuid"
)

type Photo struct {
	ID      uuid.UUID
	Content []byte
}
