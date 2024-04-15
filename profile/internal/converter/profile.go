package converter

import (
	"github.com/the-swiply/swiply-backend/profile/internal/dbmodel"
	"github.com/the-swiply/swiply-backend/profile/internal/domain"
)

func ProfileFromDBModelToDomain(interests []dbmodel.Interest, profile dbmodel.Profile) (domain.Profile, error) {
	var intrs []domain.Interest
	for _, intr := range interests {
		intrs = append(intrs, InterestFromDBModelToDomain(intr))
	}

	p := domain.Profile{
		ID:        profile.ID,
		Email:     profile.Email,
		Name:      profile.Name,
		Interests: intrs,
		BirthDay:  profile.BirthDay,
		Info:      profile.Info,
		Location: domain.Location{
			Lat:  profile.Lat,
			Long: profile.Long,
		},
	}

	if err := p.Gender.Set(profile.Gender); err != nil {
		return p, err
	}

	if err := p.Subscription.Set(profile.Subscription); err != nil {
		return p, err
	}

	return p, nil
}

func ProfileFromDomainToDBModel(profile domain.Profile) dbmodel.Profile {
	var interests []int64
	for _, interest := range profile.Interests {
		interests = append(interests, interest.ID)
	}
	return dbmodel.Profile{
		ID:           profile.ID,
		Email:        profile.Email,
		Name:         profile.Name,
		Interests:    interests,
		BirthDay:     profile.BirthDay,
		Gender:       string(profile.Gender),
		Lat:          profile.Location.Lat,
		Long:         profile.Location.Long,
		Info:         profile.Info,
		Subscription: string(profile.Subscription),
	}
}
