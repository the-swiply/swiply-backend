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
	updateInfoTable   = "update_info"
	interactionsTable = "interactions"
	statisticsTable   = "statistics"
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

func (d *DataProviderRepository) UpdateLastProfileUpdate(ctx context.Context, ts time.Time) error {
	return d.setLastUpdate(ctx, "profile", ts)
}

func (d *DataProviderRepository) UpdateLastInteractionUpdate(ctx context.Context, ts time.Time) error {
	return d.setLastUpdate(ctx, "interaction", ts)
}

func (d *DataProviderRepository) setLastUpdate(ctx context.Context, table string, ts time.Time) error {
	q := fmt.Sprintf(`UPDATE %s SET last_update = $1 WHERE entity = $2`, updateInfoTable)
	_, err := d.db.Exec(ctx, q, ts, table)
	if err != nil {
		return fmt.Errorf("can't exec query: %w", err)
	}

	return nil
}

func (d *DataProviderRepository) UpsertProfiles(ctx context.Context, profiles []domain.Profile) error {
	q := fmt.Sprintf(`INSERT INTO %s (id, updated_at)
VALUES (@id, @updated_at)
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

func (d *DataProviderRepository) CalculateRatings(ctx context.Context) (map[string]float64, error) {
	q := fmt.Sprintf(`WITH 
     positives AS (SELECT count(1) AS cnt, "to"
                   FROM %s
                   WHERE positive = true
                   GROUP BY "to"),
     alll AS (SELECT count(1) AS cnt, "to"
              FROM %s
              GROUP BY "to")
SELECT positives.cnt / alll.cnt::double precision
FROM positives JOIN alll  ON positives."to" = alll."to"`, interactionsTable, interactionsTable)

	rows, err := d.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		userID string
		rating float64
	)
	var ratings map[string]float64

	for rows.Next() {
		err = rows.Scan(&userID, &rating)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		ratings[userID] = rating
	}

	return ratings, nil
}

func (d *DataProviderRepository) UpdateStatistics(ctx context.Context, ratings map[string]float64) error {
	tx, err := d.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("can't begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		fmt.Sprintf("TRUNCATE TABLE %s", statisticsTable),
	)

	q := fmt.Sprintf(`INSERT INTO %s (user_id, like_ratio, updated_at)
VALUES (@user_id, @like_ratio, @updated_at)`, statisticsTable)

	updateTs := time.Now()
	batch := &pgx.Batch{}
	for userID, rating := range ratings {
		args := pgx.NamedArgs{
			"user_id":    userID,
			"like_ratio": rating,
			"updated_at": updateTs,
		}

		batch.Queue(q, args)
	}

	results := d.db.SendBatch(ctx, batch)
	defer results.Close()

	for range ratings {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("unable to update statistic row: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("can't commit tx: %w", err)
	}

	return nil
}
