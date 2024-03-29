package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/the-swiply/swiply-backend/pkg/houston/auf"

	"github.com/the-swiply/swiply-backend/profile/internal/converter"
	"github.com/the-swiply/swiply-backend/profile/internal/dbmodel"
	"github.com/the-swiply/swiply-backend/profile/internal/domain"
)

type PhotoRepository interface {
	List(ctx context.Context, userID uuid.UUID) ([]dbmodel.Photo, error)
	Get(ctx context.Context, userID, photoID uuid.UUID) (dbmodel.Photo, error)
	Create(ctx context.Context, userID uuid.UUID, photo dbmodel.Photo) error
	Delete(ctx context.Context, userID, photoID uuid.UUID) error
}

type ProfileRepository interface {
	ListInterests(ctx context.Context) ([]dbmodel.Interest, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (dbmodel.Profile, error)
	CreateProfile(ctx context.Context, profile dbmodel.Profile) error
	UpdateProfile(ctx context.Context, profile dbmodel.Profile) error
	CreateInteraction(ctx context.Context, interaction dbmodel.Interaction) error
	LikedProfiles(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	LikedMeProfiles(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}

type ProfileService struct {
	cfg               ProfileConfig
	photoRepository   PhotoRepository
	profileRepository ProfileRepository
}

func NewProfileService(cfg ProfileConfig, photoRepository PhotoRepository, profileRepository ProfileRepository) *ProfileService {
	return &ProfileService{
		cfg:               cfg,
		photoRepository:   photoRepository,
		profileRepository: profileRepository,
	}
}

func (p *ProfileService) Create(ctx context.Context, profile domain.Profile) error {
	userID := auf.ExtractUserIDFromContext[uuid.UUID](ctx)
	profile.ID = userID

	for _, photo := range profile.Photos {
		photo.ID = uuid.New()
		ph := converter.PhotoFromDomainToDBModel(photo)

		err := p.photoRepository.Create(ctx, userID, ph)
		if err != nil {
			return err
		}
	}

	pr := converter.ProfileFromDomainToDBModel(profile)

	return p.profileRepository.CreateProfile(ctx, pr)
}

func (p *ProfileService) Get(ctx context.Context, userID uuid.UUID) (domain.Profile, error) {
	photos, err := p.photoRepository.List(ctx, userID)
	if err != nil {
		return domain.Profile{}, err
	}

	allInterests, err := p.profileRepository.ListInterests(ctx)
	if err != nil {
		return domain.Profile{}, err
	}

	profile, err := p.profileRepository.GetProfile(ctx, userID)
	if err != nil {
		return domain.Profile{}, err
	}

	var userInterests map[int64]struct{}
	for _, interest := range profile.Interests {
		userInterests[interest] = struct{}{}
	}

	var interests []dbmodel.Interest
	for _, interest := range allInterests {
		if _, ok := userInterests[interest.ID]; ok {
			interests = append(interests, interest)
		}
	}

	return converter.ProfileFromDBModelToDomain(interests, photos, profile)
}

func (p *ProfileService) GetMyProfile(ctx context.Context) (domain.Profile, error) {
	userID := auf.ExtractUserIDFromContext[uuid.UUID](ctx)
	return p.Get(ctx, userID)
}

func (p *ProfileService) ListInterests(ctx context.Context) ([]domain.Interest, error) {
	interests, err := p.profileRepository.ListInterests(ctx)
	if err != nil {
		return nil, err
	}

	var domainInterests []domain.Interest
	for _, interest := range interests {
		domainInterests = append(domainInterests, converter.InterestFromDBModelToDomain(interest))
	}

	return domainInterests, nil
}

func (p *ProfileService) UpdateProfile(ctx context.Context, profile domain.Profile) error {
	userID := auf.ExtractUserIDFromContext[uuid.UUID](ctx)
	profile.ID = userID

	for _, photo := range profile.Photos {
		ph := converter.PhotoFromDomainToDBModel(photo)
		if ph.ID != uuid.Nil {
			ph.ID = uuid.New()
			err := p.photoRepository.Create(ctx, userID, ph)
			if err != nil {
				return err
			}
		}
	}

	pr := converter.ProfileFromDomainToDBModel(profile)

	return p.profileRepository.UpdateProfile(ctx, pr)
}
