package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	profileTable     = "profile"
	interactionTable = "interaction"
	recHistoryTable  = "recommendation_history"
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
AND mod(st.like_ratio - (SELECT like_ratio FROM statistic WHERE st.user_id = $2)) < $3
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
	q := fmt.Sprintf(`WITH excluded AS (SELECT recommendation AS id
FROM %s
WHERE user_id = $1
AND dttm >= now() - interval '%dh')

SELECT user_id FROM statistic AS st
LEFT JOIN excluded ON excluded.id = st.user_id
WHERE excluded.id IS NULL
ORDER BY random()
LIMIT %d`, recHistoryTable, excludeLastHours, limit)

	rows, err := e.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
}
