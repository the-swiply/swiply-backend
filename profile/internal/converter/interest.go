package converter

import (
	"github.com/the-swiply/swiply-backend/profile/internal/dbmodel"
	"github.com/the-swiply/swiply-backend/profile/internal/domain"
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
