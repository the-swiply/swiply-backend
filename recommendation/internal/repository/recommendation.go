package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	profileTable     = "profile"
	interactionTable = "interaction"
)

type RecommendationRepository struct {
	db *pgxpool.Pool
}

func NewRecommendationRepository(db *pgxpool.Pool) *RecommendationRepository {
	return &RecommendationRepository{
		db: db,
	}
}
