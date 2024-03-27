package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/the-swiply/swiply-backend/profile/internal/dbmodel"
)

const (
	interestTable    = "interest"
	profileTable     = "profile"
	interactionTable = "interaction"
)

type ProfileRepository struct {
	db *pgxpool.Pool
}

func NewProfileRepository(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (p *ProfileRepository) ListInterests(ctx context.Context) ([]dbmodel.Interest, error) {

}

func (p *ProfileRepository) GetProfile(ctx context.Context, userID uuid.UUID) ([]dbmodel.Profile, error) {

}

func (p *ProfileRepository) CreateProfile(ctx context.Context, profile dbmodel.Profile) error {

}

func (p *ProfileRepository) UpdateProfile(ctx context.Context, profile dbmodel.Profile) error {

}

func (p *ProfileRepository) CreateInteraction(ctx context.Context, interaction dbmodel.Interaction) error {

}

func (p *ProfileRepository) LikedProfiles(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {

}

func (p *ProfileRepository) LikedMeProfiles(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {

}
