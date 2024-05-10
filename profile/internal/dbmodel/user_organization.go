package dbmodel

import (
	"github.com/google/uuid"
)

type UserOrganization struct {
	ID             int64
	Name           string
	ProfileID      uuid.UUID `db:"profile_id"`
	OrganizationID int64     `db:"organization_id"`
	Email          string
	IsValid        bool `db:"is_valid"`
}
