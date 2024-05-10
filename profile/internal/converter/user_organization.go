package converter

import (
	"github.com/the-swiply/swiply-backend/profile/internal/dbmodel"
	"github.com/the-swiply/swiply-backend/profile/internal/domain"
)

func UserOrganizationFromDBModelToDomain(org dbmodel.UserOrganization) domain.UserOrganization {
	return domain.UserOrganization{
		ID:             org.ID,
		Name:           org.Name,
		OrganizationID: org.OrganizationID,
		Email:          org.Email,
		IsValid:        org.IsValid,
	}
}
