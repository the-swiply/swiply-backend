package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/the-swiply/swiply-backend/profile/internal/domain"
)

const (
	photoTable = "photo"
)

type PhotoRepository struct {
	db *pgxpool.Pool
}

func NewPhotoRepository(db *pgxpool.Pool) *PhotoRepository {
	return &PhotoRepository{db: db}
}

func (p *PhotoRepository) Reorder(ctx context.Context, userID uuid.UUID, photoIDs []uuid.UUID) error {
	q := fmt.Sprintf(`UPDATE %s
SET photo_ids = $1
WHERE id = $2`, photoTable)

	_, err := p.db.Exec(ctx, q, photoIDs, userID)
	return err
}

func (p *PhotoRepository) List(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	q := fmt.Sprintf(`SELECT photo_ids FROM %s WHERE id = $1`, photoTable)

	row := p.db.QueryRow(ctx, q, userID)

	var photoIDs []uuid.UUID
	err := row.Scan(&photoIDs)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrEntityIsNotExists
	}
	if err != nil {
		return nil, fmt.Errorf("can't get photo ids: %w", err)
	}

	return photoIDs, nil
}

func (p *PhotoRepository) Append(ctx context.Context, userID, photoID uuid.UUID) error {
	q := fmt.Sprintf(`UPDATE %s
SET photo_ids = array_append(photo_ids, $1)
WHERE id = $2`, photoTable)

	_, err := p.db.Exec(ctx, q, photoID, userID)
	return err
}

func (p *PhotoRepository) Delete(ctx context.Context, userID, photoID uuid.UUID) error {
	q := fmt.Sprintf(`UPDATE %s
SET photo_ids = array_remove(photo_ids, $1)
WHERE id = $2`, photoTable)

	_, err := p.db.Exec(ctx, q, photoID, userID)
	return err
}

func (p *PhotoRepository) Create(ctx context.Context, userID uuid.UUID, photoIDs []uuid.UUID) error {
	q := fmt.Sprintf(`INSERT INTO %s (id, photo_ids)
VALUES ($1, $2)`)

	_, err := p.db.Exec(ctx, q, userID, photoIDs)
	return err
}
