package domain

type UserOrganization struct {
	ID             int64
	Name           string
	OrganizationID int64
	Email          string
	IsValid        bool
}
