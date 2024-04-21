package domain

import (
	"errors"

	"github.com/google/uuid"
)

type interactionType string

func (i *interactionType) Set(value string) error {
	switch value {
	case string(InteractionTypeUnspecified):
		*i = InteractionTypeUnspecified
	case string(InteractionTypeLike):
		*i = InteractionTypeLike
	case string(InteractionTypeDislike):
		*i = InteractionTypeDislike
	}

	return errors.New("unknown interaction")
}

const (
	InteractionTypeUnspecified interactionType = ""
	InteractionTypeLike        interactionType = "LIKE"
	InteractionTypeDislike     interactionType = "DISLIKE"
)

type Interaction struct {
	ID   int64
	From uuid.UUID
	To   uuid.UUID
	Type interactionType
}
