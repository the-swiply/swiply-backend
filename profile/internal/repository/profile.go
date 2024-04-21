package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/the-swiply/swiply-backend/profile/internal/dbmodel"
	"github.com/the-swiply/swiply-backend/profile/internal/domain"
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
	q := fmt.Sprintf(`SELECT id, definition FROM %s`, interestTable)

	rows, err := p.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[dbmodel.Interest])
}

func (p *ProfileRepository) GetProfile(ctx context.Context, userID uuid.UUID) (dbmodel.Profile, error) {
	q := fmt.Sprintf(`SELECT id, email, "name", interests, birth_day, gender, info, subscription, location_lat, location_long, updated_at FROM %s
WHERE id = $1
LIMIT 1`, profileTable)

	rows, err := p.db.Query(ctx, q, userID)
	if err != nil {
		return dbmodel.Profile{}, err
	}
	defer rows.Close()

	profiles, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbmodel.Profile])
	if err != nil {
		return dbmodel.Profile{}, err
	}

	if len(profiles) == 0 {
		return dbmodel.Profile{}, domain.ErrEntityIsNotExists
	}

	return profiles[0], nil
}

func (p *ProfileRepository) CreateProfile(ctx context.Context, profile dbmodel.Profile) error {
	q := fmt.Sprintf(`INSERT INTO %s (id, email, "name", interests, birth_day, gender, info, subscription, location_lat, location_long)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, profileTable)

	_, err := p.db.Exec(ctx, q, profile.ID, profile.Email, profile.Name, profile.Interests, profile.BirthDay,
		profile.Gender, profile.Info, profile.Subscription, profile.Lat, profile.Long)
	return err
}

func (p *ProfileRepository) UpdateProfile(ctx context.Context, profile dbmodel.Profile) error {
	q := fmt.Sprintf(`UPDATE %s
SET "name"        = $1
    interests     = $2
    birth_day     = $3
    gender        = $4
    info          = $5
    subscription  = $6
    location_lat  = $7
	location_long = $8
	updated_at    = $9
WHERE id = $10`, profileTable)

	_, err := p.db.Exec(ctx, q, profile.Name, profile.BirthDay, profile.Gender, profile.Interests,
		profile.Subscription, profile.Lat, profile.Long, profile.UpdatedAt, profile.ID)
	return err
}

func (p *ProfileRepository) CreateInteraction(ctx context.Context, interaction dbmodel.Interaction) (int64, error) {
	q := fmt.Sprintf(`INSERT INTO %s ("from", "to", "type", created_at)
VALUES ($1, $2, $3, $4) 
RETURNING id`, interactionTable)

	var id int64
	row := p.db.QueryRow(ctx, q, interaction.From, interaction.To, interaction.Type, interaction.CreatedAt)

	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (p *ProfileRepository) LikedProfiles(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	q := fmt.Sprintf(`SELECT "to" FROM %s
WHERE "from" = $1 AND "type" = $2`, interactionTable)

	row := p.db.QueryRow(ctx, q, userID)

	var users []uuid.UUID
	err := row.Scan(&users)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("can't get liked profiles: %w", err)
	}

	return users, nil
}

func (p *ProfileRepository) LikedMeProfiles(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	q := fmt.Sprintf(`SELECT "from" FROM %s
WHERE "to" = $1 AND "type" = $2`, interactionTable)

	row := p.db.QueryRow(ctx, q, userID)

	var users []uuid.UUID
	err := row.Scan(&users)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("can't get liked me profiles: %w", err)
	}

	return users, nil
}

func (p *ProfileRepository) ListInteractions(ctx context.Context, createdAt time.Time) ([]dbmodel.Interaction, error) {
	q := fmt.Sprintf(`SELECT (id, "from", "to", "type", created_at) FROM %s
WHERE $1 < created_at`, interactionTable)

	rows, err := p.db.Query(ctx, q, createdAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[dbmodel.Interaction])
}

func (p *ProfileRepository) ListProfiles(ctx context.Context, updatedAt time.Time) ([]dbmodel.Profile, error) {
	q := fmt.Sprintf(`SELECT id, email, "name", interests, birth_day, gender, info, subscription, location_lat, location_long, updated_at FROM %s
WHERE $1 < updated_at`, profileTable)

	rows, err := p.db.Query(ctx, q, updatedAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[dbmodel.Profile])
}
