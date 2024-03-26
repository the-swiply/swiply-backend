package converter

import (
	"github.com/the-swiply/swiply-backend/profile/internal/dbmodel"
	"github.com/the-swiply/swiply-backend/profile/internal/domain"
)

func PhotoFromDBModelToDomain(photo dbmodel.Photo) domain.Photo {
	return domain.Photo{
		ID:      photo.ID,
		Content: photo.Content,
	}
}

func PhotoFromDomainToDBModel(photo domain.Photo) dbmodel.Photo {
	return dbmodel.Photo{
		ID:      photo.ID,
		Content: photo.Content,
	}
}
