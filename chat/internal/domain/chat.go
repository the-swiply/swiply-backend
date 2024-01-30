package domain

import "github.com/google/uuid"

type Chat struct {
	ID      int64
	Members []uuid.UUID
}
