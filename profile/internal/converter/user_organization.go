package converter

import (
	"github.com/the-swiply/swiply-backend/profile/internal/dbmodel"
	"github.com/the-swiply/swiply-backend/profile/internal/domain"
	"github.com/the-swiply/swiply-backend/profile/pkg/api/profile"
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

func UserOrganizationFromDomainToProto(org domain.UserOrganization) *profile.UserOrganization {
	return &profile.UserOrganization{
		Id:             org.ID,
		Name:           org.Name,
		Email:          org.Email,
		IsValid:        org.IsValid,
		OrganizationId: org.OrganizationID,
	}
}

func UserOrganizationFromProtoToDomain(org *profile.UserOrganization) domain.UserOrganization {
	return domain.UserOrganization{
		ID:             org.GetId(),
		Name:           org.GetName(),
		OrganizationID: org.GetOrganizationId(),
		Email:          org.GetEmail(),
		IsValid:        org.GetIsValid(),
	}
}
