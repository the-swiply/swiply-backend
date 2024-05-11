package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/the-swiply/swiply-backend/profile/internal/converter"
	"github.com/the-swiply/swiply-backend/profile/internal/dbmodel"
	"github.com/the-swiply/swiply-backend/profile/internal/domain"
)

type ProfileRepository interface {
	ListInterests(ctx context.Context) ([]dbmodel.Interest, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (dbmodel.Profile, error)
	CreateProfile(ctx context.Context, profile dbmodel.Profile) error
	UpdateProfile(ctx context.Context, profile dbmodel.Profile) error
	CreateInteraction(ctx context.Context, interaction dbmodel.Interaction) (int64, error)
	LikedProfiles(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	LikedMeProfiles(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	ListInteractions(ctx context.Context, createdAt time.Time) ([]dbmodel.Interaction, error)
	ListProfiles(ctx context.Context, updatedAt time.Time) ([]dbmodel.Profile, error)
	ChangeAvailability(ctx context.Context, isBlocked bool, userID uuid.UUID) error
	CreateUserOrganization(ctx context.Context, profileID uuid.UUID, email string) (dbmodel.UserOrganization, error)
	DeleteUserOrganization(ctx context.Context, userID uuid.UUID, id int64) error
	ListUserOrganizations(ctx context.Context, userID uuid.UUID) ([]dbmodel.UserOrganization, error)
	ValidateUserOrganization(ctx context.Context, userID uuid.UUID, id int64) error
	ListMatches(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}

type ProfileService struct {
	cfg  ProfileConfig
	repo ProfileRepository
}

func NewProfileService(cfg ProfileConfig, profileRepository ProfileRepository) *ProfileService {
	return &ProfileService{
		cfg:  cfg,
		repo: profileRepository,
	}
}

func (p *ProfileService) Create(ctx context.Context, profile domain.Profile) error {
	return p.repo.CreateProfile(ctx, converter.ProfileFromDomainToDBModel(profile))
}

func (p *ProfileService) Update(ctx context.Context, profile domain.Profile) error {
	return p.repo.UpdateProfile(ctx, converter.ProfileFromDomainToDBModel(profile))
}

func (p *ProfileService) Get(ctx context.Context, userID uuid.UUID) (domain.Profile, error) {
	profile, err := p.repo.GetProfile(ctx, userID)
	if err != nil {
		return domain.Profile{}, err
	}

	allInterests, err := p.repo.ListInterests(ctx)
	if err != nil {
		return domain.Profile{}, err
	}

	organizations, err := p.repo.ListUserOrganizations(ctx, userID)
	if err != nil {
		return domain.Profile{}, err
	}

	mp := make(map[int64]dbmodel.Interest)
	for _, interest := range allInterests {
		mp[interest.ID] = interest
	}

	var userInterests []dbmodel.Interest
	for _, interest := range profile.Interests {
		userInterests = append(userInterests, mp[interest])
	}

	return converter.ProfileFromDBModelToDomain(userInterests, organizations, profile)
}

func (p *ProfileService) CreateInteraction(ctx context.Context, interaction domain.Interaction) error {
	_, err := p.repo.CreateInteraction(ctx, converter.InteractionFromDomainToDBModel(interaction))
	return err
}

func (p *ProfileService) GetLikedProfiles(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	return p.repo.LikedProfiles(ctx, userID)
}

func (p *ProfileService) GetLikedMeProfiles(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	return p.repo.LikedMeProfiles(ctx, userID)
}

func (p *ProfileService) ListInterests(ctx context.Context) ([]domain.Interest, error) {
	interests, err := p.repo.ListInterests(ctx)
	if err != nil {
		return nil, err
	}

	var inters []domain.Interest

	for _, interest := range interests {
		inters = append(inters, converter.InterestFromDBModelToDomain(interest))
	}

	return inters, nil
}

func (p *ProfileService) ListInteractions(ctx context.Context, createdAt time.Time) ([]domain.Interaction, error) {
	interactions, err := p.repo.ListInteractions(ctx, createdAt)
	if err != nil {
		return nil, err
	}

	var interacts []domain.Interaction

	for _, interaction := range interactions {
		interact, err := converter.InteractionFromDBModelToDomain(interaction)
		if err != nil {
			return nil, err
		}
		interacts = append(interacts, interact)
	}

	return interacts, nil
}

func (p *ProfileService) ListProfiles(ctx context.Context, updatedAt time.Time) ([]domain.Profile, error) {
	profiles, err := p.repo.ListProfiles(ctx, updatedAt)
	if err != nil {
		return nil, err
	}

	allInterests, err := p.repo.ListInterests(ctx)
	if err != nil {
		return nil, err
	}

	mp := make(map[int64]dbmodel.Interest)
	for _, interest := range allInterests {
		mp[interest.ID] = interest
	}

	var profs []domain.Profile

	for _, profile := range profiles {
		var userInterests []dbmodel.Interest
		for _, interest := range profile.Interests {
			userInterests = append(userInterests, mp[interest])
		}

		organizations, err := p.repo.ListUserOrganizations(ctx, profile.ID)
		if err != nil {
			continue
		}

		if pr, err := converter.ProfileFromDBModelToDomain(userInterests, organizations, profile); err == nil {
			profs = append(profs, pr)
		}
	}

	return profs, nil
}

func (p *ProfileService) ChangeAvailability(ctx context.Context, isBlocked bool, userID uuid.UUID) error {
	return p.repo.ChangeAvailability(ctx, isBlocked, userID)
}

func (p *ProfileService) AddUserOrganization(ctx context.Context, userID uuid.UUID, email string) (domain.UserOrganization, error) {
	org, err := p.repo.CreateUserOrganization(ctx, userID, email)
	if err != nil {
		return domain.UserOrganization{}, err
	}

	return converter.UserOrganizationFromDBModelToDomain(org), nil
}

func (p *ProfileService) DeleteUserOrganization(ctx context.Context, userID uuid.UUID, id int64) error {
	return p.repo.DeleteUserOrganization(ctx, userID, id)
}

func (p *ProfileService) SendAuthorizationCode(ctx context.Context, userID uuid.UUID, email string) error {
	return nil
}

func (p *ProfileService) ValidateUserOrganization(ctx context.Context, userID uuid.UUID, id int64, code string) error {
	return p.repo.ValidateUserOrganization(ctx, userID, id)
}

func (p *ProfileService) ListMatches(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	return p.repo.ListMatches(ctx, userID)
}
