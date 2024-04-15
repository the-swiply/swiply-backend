package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/the-swiply/swiply-backend/profile/internal/dbmodel"
)

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
	profileRepository ProfileRepository
}

func NewProfileService(cfg ProfileConfig, profileRepository ProfileRepository) *ProfileService {
	return &ProfileService{
		cfg:               cfg,
		profileRepository: profileRepository,
	}
}

//func (p *ProfileService) Create(ctx context.Context, profile domain.Profile) error {
//	userID := auf.ExtractUserIDFromContext[uuid.UUID](ctx)
//	profile.ID = userID
//}
//
//func (p *ProfileService) Get(ctx context.Context, userID uuid.UUID) (domain.Profile, error) {
//
//}
//
//func (p *ProfileService) GetMyProfile(ctx context.Context) (domain.Profile, error) {
//	userID := auf.ExtractUserIDFromContext[uuid.UUID](ctx)
//	return p.Get(ctx, userID)
//}
//
//func (p *ProfileService) ListInterests(ctx context.Context) ([]domain.Interest, error) {
//
//}
//
//func (p *ProfileService) UpdateProfile(ctx context.Context, profile domain.Profile) error {
//
//}
