package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/the-swiply/swiply-backend/pkg/houston/auf"
	"github.com/the-swiply/swiply-backend/recommendation/internal/domain"
)

type RecommendationRepository interface {
}

type RecommendationService struct {
	cfg     RecommendationConfig
	recRepo RecommendationRepository
}

func NewRecommendationService(cfg RecommendationConfig, recRepo RecommendationRepository) *RecommendationService {
	return &RecommendationService{
		cfg:     cfg,
		recRepo: recRepo,
	}
}

func (r *RecommendationService) Recommend(ctx context.Context, count int64) ([]domain.Recommendation, error) {
	userID := auf.ExtractUserIDFromContext[uuid.UUID](ctx)
	_ = userID
	recs := make([]domain.Recommendation, count)
	return recs, nil
}
