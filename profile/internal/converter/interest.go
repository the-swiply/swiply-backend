package converter

import (
	"github.com/the-swiply/swiply-backend/profile/internal/dbmodel"
	"github.com/the-swiply/swiply-backend/profile/internal/domain"
	"github.com/the-swiply/swiply-backend/profile/pkg/api/profile"
)

func InterestFromDBModelToDomain(interest dbmodel.Interest) domain.Interest {
	return domain.Interest{
		ID:         interest.ID,
		Definition: interest.Definition,
	}
}

func InterestFromDomainToDBModel(interest domain.Interest) dbmodel.Interest {
	return dbmodel.Interest{
		ID:         interest.ID,
		Definition: interest.Definition,
	}
}

func InterestFromDomainToProto(interest domain.Interest) *profile.Interest {
	return &profile.Interest{
		Id:         interest.ID,
		Definition: interest.Definition,
	}
}

func InterestFromProtoToDomain(interest *profile.Interest) domain.Interest {
	return domain.Interest{
		ID:         interest.GetId(),
		Definition: interest.GetDefinition(),
	}
}
