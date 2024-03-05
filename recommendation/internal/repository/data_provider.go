package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/the-swiply/swiply-backend/recommendation/internal/domain"
	"time"
)

const (
	updateInfoTable = "update_info"
)

type DataProviderRepository struct {
	db *pgxpool.Pool
}

func NewDataProviderRepository(db *pgxpool.Pool) *DataProviderRepository {
	return &DataProviderRepository{
		db: db,
	}
}
func (d *DataProviderRepository) GetLastProfileUpdate(ctx context.Context) (time.Time, error) {
	return d.getLastUpdate(ctx, "profile")
}

func (d *DataProviderRepository) GetLastInteractionUpdate(ctx context.Context) (time.Time, error) {
	return d.getLastUpdate(ctx, "interaction")
}

func (d *DataProviderRepository) getLastUpdate(ctx context.Context, table string) (time.Time, error) {
	q := fmt.Sprintf(`SELECT last_update FROM %s WHERE entity = $1`, updateInfoTable)
	row := d.db.QueryRow(ctx, q, table)

	var lastUpdate time.Time
	err := row.Scan(&lastUpdate)
	if err != nil {
		return time.Time{}, err
	}

	return lastUpdate, nil
}

func (d *DataProviderRepository) UpsertProfiles(ctx context.Context, profiles []domain.Profile) error {
	q := fmt.Sprintf(`INSERT INTO %s (id, from, to, positive, updated_at)
VALUES (@id, @from, @to, @positive, @updated_at)
ON CONFLICT(id) 
DO UPDATE SET
id = @id,
from = @from,
to = @to,
positive = @positive,
updated_at = @updated_at`, profileTable)

	batch := &pgx.Batch{}
	for _, profile := range profiles {
		args := pgx.NamedArgs{
			"id":         profile.ID,
			"updated_at": profile.UpdatedAt,
		}

		batch.Queue(q, args)
	}

	results := d.db.SendBatch(ctx, batch)
	defer results.Close()

	for range profiles {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("unable to upsert profile row: %w", err)
		}
	}

	return results.Close()
}

func (d *DataProviderRepository) UpsertInteractions(ctx context.Context, interactions []domain.Interaction) error {
	q := fmt.Sprintf(`INSERT INTO %s (id, from, to, positive, updated_at)
VALUES (@id, @from, @to, @positive, @updated_at)
ON CONFLICT(id) 
DO UPDATE SET
id = @id,
from = @from,
to = @to,
positive = @positive,
updated_at = @updated_at`, interactionTable)

	batch := &pgx.Batch{}
	for _, interaction := range interactions {
		args := pgx.NamedArgs{
			"id":         interaction.ID,
			"from":       interaction.From,
			"to":         interaction.To,
			"positive":   interaction.Positive,
			"updated_at": interaction.UpdatedAt,
		}

		batch.Queue(q, args)
	}

	results := d.db.SendBatch(ctx, batch)
	defer results.Close()

	for range interactions {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("unable to upsert interaction row: %w", err)
		}
	}

	return results.Close()
}
