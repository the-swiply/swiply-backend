package converter

import (
	"github.com/google/uuid"

	"github.com/the-swiply/swiply-backend/profile/internal/dbmodel"
	"github.com/the-swiply/swiply-backend/profile/internal/domain"
	"github.com/the-swiply/swiply-backend/profile/pkg/api/profile"
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

func PhotoFromProtoToDomain(photo *profile.ProfilePhoto) domain.Photo {
	return domain.Photo{
		ID:      uuid.MustParse(photo.GetId()),
		Content: photo.GetContent(),
	}
}

func PhotoFromDomainToProto(photo domain.Photo) *profile.ProfilePhoto {
	return &profile.ProfilePhoto{
		Id:      photo.ID.String(),
		Content: photo.Content,
	}
}
