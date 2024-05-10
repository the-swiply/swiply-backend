package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/the-swiply/swiply-backend/profile/internal/dbmodel"
	"github.com/the-swiply/swiply-backend/profile/internal/domain"
)

const (
	interestTable         = "interest"
	profileTable          = "profile"
	interactionTable      = "interaction"
	organizationTable     = "organization"
	userOrganizationTable = "user_organization"
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
	q := fmt.Sprintf(`SELECT id, email, "name", city, "work", education, is_blocked, interests, birth_day, gender, info, subscription, location_lat, location_long, updated_at FROM %s
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
	q := fmt.Sprintf(`INSERT INTO %s (id, email, "name", city, "work", education, is_blocked, interests, birth_day, gender, info, subscription, location_lat, location_long)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`, profileTable)

	_, err := p.db.Exec(ctx, q, profile.ID, profile.Email, profile.Name, profile.City, profile.Work, profile.Education,
		profile.IsBlocked, profile.Interests, profile.BirthDay, profile.Gender, profile.Info, profile.Subscription,
		profile.Lat, profile.Long)
	return err
}

func (p *ProfileRepository) UpdateProfile(ctx context.Context, profile dbmodel.Profile) error {
	q := fmt.Sprintf(`UPDATE %s
SET "name"        = $1
	city		  = $2
	"work" 		  = $3 
	education	  = $4 
    interests     = $5
    birth_day     = $6
    gender        = $7
    info          = $8
    subscription  = $9
    location_lat  = $10
	location_long = $11
	updated_at    = $12
WHERE id = $13`, profileTable)

	_, err := p.db.Exec(ctx, q, profile.Name, profile.City, profile.Work, profile.Education, profile.BirthDay, profile.Gender, profile.Interests,
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

	row := p.db.QueryRow(ctx, q, userID, domain.InteractionTypeLike)

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
WHERE "to" = $1 AND "type" = $2 AND "from" NOT IN (SELECT "to" FROM %s WHERE "from" = $1)`, interactionTable, interactionTable)

	row := p.db.QueryRow(ctx, q, userID, domain.InteractionTypeLike)

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
	q := fmt.Sprintf(`SELECT id, "from", "to", "type", created_at FROM %s
WHERE $1 < created_at`, interactionTable)

	rows, err := p.db.Query(ctx, q, createdAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[dbmodel.Interaction])
}

func (p *ProfileRepository) ListProfiles(ctx context.Context, updatedAt time.Time) ([]dbmodel.Profile, error) {
	q := fmt.Sprintf(`SELECT id, email, "name", city, "work", education, is_blocked, interests, birth_day, gender, info, subscription, location_lat, location_long, updated_at FROM %s
WHERE $1 < updated_at`, profileTable)

	rows, err := p.db.Query(ctx, q, updatedAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[dbmodel.Profile])
}

func (p *ProfileRepository) ChangeAvailability(ctx context.Context, isBlocked bool, userID uuid.UUID) error {
	q := fmt.Sprintf(`UPDATE %s
SET is_blocked = $1
WHERE id = $2`, profileTable)

	_, err := p.db.Exec(ctx, q, isBlocked, userID)
	return err
}

func (p *ProfileRepository) ListOrganizations(ctx context.Context) ([]dbmodel.Organization, error) {
	q := fmt.Sprintf(`SELECT id, "name", pattern FROM %s`, organizationTable)
	rows, err := p.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[dbmodel.Organization])
}

func (p *ProfileRepository) CreateUserOrganization(ctx context.Context, profileID uuid.UUID, email string) (dbmodel.UserOrganization, error) {
	var (
		userOrg dbmodel.UserOrganization
		exist   bool
	)

	organizations, err := p.ListOrganizations(ctx)
	if err != nil {
		return dbmodel.UserOrganization{}, err
	}

	for _, org := range organizations {
		if strings.Split(email, "@")[1] == org.Pattern {
			exist = true
			userOrg.Name = org.Name
			userOrg.ProfileID = profileID
			userOrg.OrganizationID = org.ID
			userOrg.Email = email
			break
		}
	}

	if !exist {
		return dbmodel.UserOrganization{}, domain.ErrEntityIsNotExists
	}

	q := fmt.Sprintf(`INSERT INTO %s (profile_id, "name", organization_id, email, is_valid)
VALUES ($1, $2, $3, $4, $5) 
RETURNING id`, userOrganizationTable)

	var id int64
	row := p.db.QueryRow(ctx, q, userOrg.ProfileID, userOrg.Name, userOrg.OrganizationID, userOrg.Email, userOrg.IsValid)

	err = row.Scan(&id)
	if err != nil {
		return dbmodel.UserOrganization{}, err
	}

	userOrg.ID = id
	return userOrg, nil
}

func (p *ProfileRepository) DeleteUserOrganization(ctx context.Context, userID uuid.UUID, id int64) error {
	q := fmt.Sprintf(`DELETE FROM %s
WHERE id = $1 AND profile_id = $2`, userOrganizationTable)

	_, err := p.db.Exec(ctx, q, id, userID, id)
	if err != nil {
		return fmt.Errorf("can't delete user organization in db: %w", err)
	}

	return nil
}

func (p *ProfileRepository) ListUserOrganizations(ctx context.Context, userID uuid.UUID) ([]dbmodel.UserOrganization, error) {
	q := fmt.Sprintf(`SELECT id, profile_id, "name", organization_id, email, is_valid FROM %s
WHERE profile_id = $1`, profileTable)

	rows, err := p.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[dbmodel.UserOrganization])
}

func (p *ProfileRepository) ValidateUserOrganization(ctx context.Context, userID uuid.UUID, id int64) error {
	q := fmt.Sprintf(`UPDATE %s
SET is_valid = true
WHERE id = $1 and profile_id = $2`, userOrganizationTable)

	_, err := p.db.Exec(ctx, q, id, userID)
	return err
}

func (p *ProfileRepository) ListMatches(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	q := fmt.Sprintf(`SELECT "to" FROM %s
WHERE "from" = $1 AND "type" = $2 "to" IN (SELECT "from" FROM %s
WHERE "to" = $1 AND "type" = $2)`, interactionTable, interactionTable)

	row := p.db.QueryRow(ctx, q, userID, domain.InteractionTypeLike)

	var users []uuid.UUID
	err := row.Scan(&users)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrEntityIsNotExists
	}
	if err != nil {
		return nil, fmt.Errorf("can't get match users: %w", err)
	}

	return users, nil
}
