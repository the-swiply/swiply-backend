package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/the-swiply/swiply-backend/pkg/houston/dobby"
	"github.com/the-swiply/swiply-backend/recommendation/internal/domain"
	"time"
)

type DataProviderRepository struct {
	db *pgxpool.Pool
}

func (d *DataProviderRepository) executor(ctx context.Context) dobby.Executor {
	tx := dobby.ExtractPGXTx(ctx)
	if tx != nil {
		return tx
	}

	return d.db
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
	row := d.executor(ctx).QueryRow(ctx, q, table)

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
	_, err := d.executor(ctx).Exec(ctx, q, ts, table)
	if err != nil {
		return fmt.Errorf("can't exec query: %w", err)
	}

	return nil
}

func (d *DataProviderRepository) UpsertProfiles(ctx context.Context, profiles []domain.Profile) error {
	q := fmt.Sprintf(`INSERT INTO %s (id, interests, birthday, gender, info,
                subscription_type, location_lat, location_lon, updated_at)
VALUES (@id, @interests, @birthday, @gender, @info,
                @subscription_type, @location_lat, @location_lon, @updated_at)
ON CONFLICT(id) 
DO UPDATE SET
id = @id,
interests = @interests,
birthday = @birthday,
gender = @gender,
info = @info,
subscription_type = @subscription_type,
location_lat = @location_lat,
location_lon = @location_lon,
updated_at = @updated_at`, profileTable)

	now := time.Now()
	batch := &pgx.Batch{}
	for _, profile := range profiles {
		args := pgx.NamedArgs{
			"id":                profile.ID,
			"interests":         profile.Interests,
			"birthday":          profile.BirthDay,
			"gender":            profile.Gender,
			"info":              profile.Info,
			"subscription_type": profile.SubscriptionType,
			"location_lat":      profile.LocationLat,
			"location_lon":      profile.LocationLon,
			"updated_at":        now,
		}

		batch.Queue(q, args)
	}

	results := d.executor(ctx).SendBatch(ctx, batch)
	defer results.Close()

	for range profiles {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("can't upsert profile row: %w", err)
		}
	}

	return results.Close()
}

func (d *DataProviderRepository) AddInteractions(ctx context.Context, interactions []domain.Interaction) error {
	if len(interactions) == 0 {
		return nil
	}

	now := time.Now()
	q := sq.Insert(interactionTable).Columns("id", "\"from\"", "\"to\"", "positive", "updated_at")
	for _, interaction := range interactions {
		q = q.Values(uuid.New(), interaction.From, interaction.To, interaction.Positive, now)
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("can't prepare sql: %w", err)
	}

	_, err = d.executor(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("can't execute insert: %w", err)
	}

	return nil
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
SELECT alll."to", coalesce(positives.cnt, 0) / alll.cnt::double precision
FROM positives RIGHT JOIN alll ON positives."to" = alll."to"`, interactionTable, interactionTable)

	rows, err := d.executor(ctx).Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		userID string
		rating float64
	)
	ratings := make(map[string]float64, 100)

	for rows.Next() {
		err = rows.Scan(&userID, &rating)
		if err != nil {
			return nil, fmt.Errorf("can't scan row: %w", err)
		}
		ratings[userID] = rating
	}

	return ratings, nil
}

func (d *DataProviderRepository) UpdateStatistics(ctx context.Context, ratings map[string]float64) error {
	_, err := d.executor(ctx).Exec(
		ctx,
		fmt.Sprintf("TRUNCATE TABLE %s", statisticsTable),
	)
	if err != nil {
		return err
	}

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

	results := d.executor(ctx).SendBatch(ctx, batch)
	defer results.Close()

	for range ratings {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("can't update statistic row: %w", err)
		}
	}

	return nil
}
