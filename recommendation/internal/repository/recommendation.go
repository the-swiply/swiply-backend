package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RecommendationRepository struct {
	db *pgxpool.Pool
}

func NewRecommendationRepository(db *pgxpool.Pool) *RecommendationRepository {
	return &RecommendationRepository{
		db: db,
	}
}

func (e *RecommendationRepository) GetGetRecommendationsByLikes(ctx context.Context, userID uuid.UUID, limit int64) ([]uuid.UUID, error) {
	q := fmt.Sprintf(`WITH excluded AS (SELECT "to" AS id
FROM %s
WHERE "from" = $1
AND positive = false)

SELECT "from"
FROM %s
LEFT JOIN excluded ON excluded.id = %s."from"
WHERE excluded.id IS NULL
AND "to" = $2
AND positive = true
ORDER BY random()
LIMIT %d`, interactionTable, interactionTable, interactionTable, limit)

	rows, err := e.db.Query(ctx, q, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
}

func (e *RecommendationRepository) GetRecommendationsByRating(ctx context.Context, userID uuid.UUID,
	maxRatingDelta float64, excludeLastHours int64, limit int64) ([]uuid.UUID, error) {
	q := fmt.Sprintf(`WITH excluded AS (SELECT recommendation AS id
FROM %s
WHERE user_id = $1
AND dttm >= now() - interval '%dh')

SELECT user_id FROM statistic AS st
LEFT JOIN excluded ON excluded.id = st.user_id
WHERE excluded.id IS NULL
AND abs(st.like_ratio - (SELECT like_ratio FROM statistic WHERE st.user_id = $2)) < $3
AND st.user_id != $2
ORDER BY random()
LIMIT %d`, recHistoryTable, excludeLastHours, limit)

	rows, err := e.db.Query(ctx, q, userID, userID, maxRatingDelta)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
}

func (e *RecommendationRepository) GetRecommendationsByOracle(ctx context.Context, userID uuid.UUID, excludeLastHours int64, limit int64) ([]uuid.UUID, error) {
	q := fmt.Sprintf(`WITH excluded AS (SELECT recommendation AS id
FROM %s
WHERE user_id = $1
AND dttm >= now() - interval '%dh')

SELECT user_id
FROM oracle_prediction AS pred
LEFT JOIN excluded ON excluded.id = pred.user_id
WHERE excluded.id IS NULL 
AND user_id = $1
AND recommendation != $1
ORDER BY score
LIMIT %d`, recHistoryTable, excludeLastHours, limit)

	rows, err := e.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
}

func (e *RecommendationRepository) GetRecommendationsByRandom(ctx context.Context, userID uuid.UUID, excludeLastHours int64, limit int64) ([]uuid.UUID, error) {
	q := fmt.Sprintf(`WITH excluded AS (SELECT DISTINCT recommendation AS id
FROM %s
WHERE user_id = $1
AND dttm >= now() - interval '%dh')

SELECT p.id FROM profile AS p
LEFT JOIN excluded ON excluded.id = p.id
WHERE excluded.id IS NULL
AND p.id != $1
ORDER BY random()
LIMIT %d`, recHistoryTable, excludeLastHours, limit)

	rows, err := e.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
}

func (e *RecommendationRepository) SaveUserRecommendationHistory(ctx context.Context, userID uuid.UUID, recommended []uuid.UUID) error {
	if len(recommended) == 0 {
		return nil
	}

	now := time.Now()

	q := sq.Insert(recHistoryTable).Columns("id", "user_id", "recommendation", "dttm")
	for _, recommendation := range recommended {
		q = q.Values(uuid.New(), userID, recommendation, now)
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("can't prepare sql: %w", err)
	}

	_, err = e.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("can't execute insert: %w", err)
	}

	return nil
}
