package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type gender string

func (g *gender) Set(value string) error {
	switch value {
	case string(GenderUnspecified):
		*g = GenderUnspecified
	case string(GenderMale):
		*g = GenderMale
	case string(GenderFemale):
		*g = GenderFemale
	default:
		return errors.New("unknown gender")
	}

	return nil
}

const (
	GenderUnspecified gender = ""
	GenderMale        gender = "MALE"
	GenderFemale      gender = "FEMALE"
)

type subscriptionType string

const (
	SubscriptionTypeUnspecified subscriptionType = ""
	SubscriptionTypeStandard    subscriptionType = "STANDARD"
	SubscriptionTypePrimary     subscriptionType = "PRIMARY"
)

func (s *subscriptionType) Set(value string) error {
	switch value {
	case string(SubscriptionTypeUnspecified):
		*s = SubscriptionTypeUnspecified
	case string(SubscriptionTypeStandard):
		*s = SubscriptionTypeStandard
	case string(SubscriptionTypePrimary):
		*s = SubscriptionTypePrimary
	default:
		return errors.New("unknown subscription type")
	}

	return nil
}

type Location struct {
	Lat  float64
	Long float64
}

type Profile struct {
	ID           uuid.UUID
	Email        string
	Name         string
	Interests    []Interest
	BirthDay     time.Time
	Gender       gender
	Info         string
	Subscription subscriptionType
	Location     Location
}
