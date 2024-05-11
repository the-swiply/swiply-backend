package converter

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/the-swiply/swiply-backend/profile/internal/dbmodel"
	"github.com/the-swiply/swiply-backend/profile/internal/domain"
	"github.com/the-swiply/swiply-backend/profile/pkg/api/profile"
)

func ProfileFromDBModelToDomain(interests []dbmodel.Interest, organizations []dbmodel.UserOrganization,
	profile dbmodel.Profile) (domain.Profile, error) {
	var intrs []domain.Interest
	for _, intr := range interests {
		intrs = append(intrs, InterestFromDBModelToDomain(intr))
	}

	var orgs []domain.UserOrganization
	for _, org := range organizations {
		orgs = append(orgs, UserOrganizationFromDBModelToDomain(org))
	}

	p := domain.Profile{
		ID:        profile.ID,
		Email:     profile.Email,
		Name:      profile.Name,
		City:      profile.City,
		Work:      profile.Work,
		Education: profile.Education,
		IsBlocked: profile.IsBlocked,
		Interests: intrs,
		BirthDay:  profile.BirthDay,
		Info:      profile.Info,
		Location: domain.Location{
			Lat:  profile.Lat,
			Long: profile.Long,
		},
		Organizations: orgs,
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
		City:         profile.City,
		Work:         profile.Work,
		Education:    profile.Education,
		IsBlocked:    profile.IsBlocked,
		Interests:    interests,
		BirthDay:     profile.BirthDay,
		Gender:       string(profile.Gender),
		Lat:          profile.Location.Lat,
		Long:         profile.Location.Long,
		Info:         profile.Info,
		Subscription: string(profile.Subscription),
		UpdatedAt:    time.Now().UTC(),
	}
}

func ProfileFromDomainToProto(prof domain.Profile) *profile.UserProfile {
	userProfile := &profile.UserProfile{
		Id:       prof.ID.String(),
		Email:    prof.Email,
		Name:     prof.Name,
		BirthDay: timestamppb.New(prof.BirthDay),
		Info:     prof.Info,
		Location: &profile.Location{
			Lat:  prof.Location.Lat,
			Long: prof.Location.Long,
		},
		City:      prof.City,
		Work:      prof.Work,
		Education: prof.Education,
		IsBlocked: prof.IsBlocked,
	}

	for _, interest := range prof.Interests {
		userProfile.Interests = append(userProfile.Interests, InterestFromDomainToProto(interest))
	}

	for _, org := range prof.Organizations {
		userProfile.Organizations = append(userProfile.Organizations, UserOrganizationFromDomainToProto(org))
	}

	switch prof.Gender {
	case domain.GenderUnspecified:
		userProfile.Gender = profile.Gender_GENDER_UNSPECIFIED
	case domain.GenderMale:
		userProfile.Gender = profile.Gender_MALE
	case domain.GenderFemale:
		userProfile.Gender = profile.Gender_FEMALE
	}

	switch prof.Subscription {
	case domain.SubscriptionTypeUnspecified:
		userProfile.SubscriptionType = profile.SubscriptionType_SUBSCRIPTION_TYPE_UNSPECIFIED
	case domain.SubscriptionTypeStandard:
		userProfile.SubscriptionType = profile.SubscriptionType_STANDARD
	case domain.SubscriptionTypePrimary:
		userProfile.SubscriptionType = profile.SubscriptionType_PRIMARY
	}

	return userProfile
}

func ProfileFromProtoToDomain(userProfile *profile.UserProfile) domain.Profile {
	prof := domain.Profile{
		ID:        uuid.MustParse(userProfile.GetId()),
		Email:     userProfile.GetEmail(),
		Name:      userProfile.GetName(),
		City:      userProfile.GetCity(),
		Work:      userProfile.GetWork(),
		Education: userProfile.GetEducation(),
		IsBlocked: userProfile.GetIsBlocked(),
		BirthDay:  userProfile.GetBirthDay().AsTime(),
		Info:      userProfile.GetInfo(),
		Location: domain.Location{
			Lat:  userProfile.GetLocation().GetLat(),
			Long: userProfile.GetLocation().GetLong(),
		},
	}

	for _, interest := range userProfile.GetInterests() {
		prof.Interests = append(prof.Interests, InterestFromProtoToDomain(interest))
	}

	for _, org := range userProfile.GetOrganizations() {
		prof.Organizations = append(prof.Organizations, UserOrganizationFromProtoToDomain(org))
	}

	switch userProfile.GetGender() {
	case profile.Gender_GENDER_UNSPECIFIED:
		prof.Gender = domain.GenderUnspecified
	case profile.Gender_MALE:
		prof.Gender = domain.GenderMale
	case profile.Gender_FEMALE:
		prof.Gender = domain.GenderFemale
	}

	switch userProfile.GetSubscriptionType() {
	case profile.SubscriptionType_SUBSCRIPTION_TYPE_UNSPECIFIED:
		prof.Subscription = domain.SubscriptionTypeUnspecified
	case profile.SubscriptionType_STANDARD:
		prof.Subscription = domain.SubscriptionTypeStandard
	case profile.SubscriptionType_PRIMARY:
		prof.Subscription = domain.SubscriptionTypePrimary
	}

	return prof
}
