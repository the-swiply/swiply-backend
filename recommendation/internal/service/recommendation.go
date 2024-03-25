package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/the-swiply/swiply-backend/pkg/houston/auf"
	"math"
)

const (
	defaultMaxRatingDelta = 0.2
)

type RecommendationRepository interface {
	GetGetRecommendationsByLikes(ctx context.Context, userID uuid.UUID, limit int64) ([]uuid.UUID, error)
	GetRecommendationsByRating(ctx context.Context, userID uuid.UUID, maxRatingDelta float64, excludeLastHours int64, limit int64) ([]uuid.UUID, error)
	GetRecommendationsByOracle(ctx context.Context, userID uuid.UUID, excludeLastHours int64, limit int64) ([]uuid.UUID, error)
	GetRecommendationsByRandom(ctx context.Context, userID uuid.UUID, excludeLastHours int64, limit int64) ([]uuid.UUID, error)
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

type recommendationLimits struct {
	byLikes  int64
	byRating int64
	byOracle int64
	byRandom int64
}

func (r *RecommendationService) Recommend(ctx context.Context, count int64) ([]uuid.UUID, error) {
	userID := auf.ExtractUserIDFromContext[uuid.UUID](ctx)
	recs := make([]uuid.UUID, 0, count)
	recLimits := calculateRecommendationLimits(count)

	byLikesRecs, err := r.recRepo.GetGetRecommendationsByLikes(ctx, userID, recLimits.byLikes)
	if err != nil {
		return nil, fmt.Errorf("can't get recommendations by likes: %w", err)
	}

	byRatingRecs, err := r.recRepo.GetRecommendationsByRating(ctx, userID, defaultMaxRatingDelta, r.cfg.FreezeHoursForRecommendation, recLimits.byLikes)
	if err != nil {
		return nil, fmt.Errorf("can't get recommendations by rating: %w", err)
	}

	byOracleRecs, err := r.recRepo.GetRecommendationsByOracle(ctx, userID, r.cfg.FreezeHoursForRecommendation, recLimits.byOracle)
	if err != nil {
		return nil, fmt.Errorf("can't get recommendations by oracle: %w", err)
	}

	byRandomRecs, err := r.recRepo.GetRecommendationsByRandom(ctx, userID, r.cfg.FreezeHoursForRecommendation, recLimits.byRandom)
	if err != nil {
		return nil, fmt.Errorf("can't get recommendations by random: %w", err)
	}

	recs = append(recs, byLikesRecs...)
	recs = append(recs, byRatingRecs...)
	recs = append(recs, byOracleRecs...)
	recs = append(recs, byRandomRecs...)

	return recs, nil
}

func calculateRecommendationLimits(totalLimit int64) recommendationLimits {
	limits := recommendationLimits{}

	// Default rule is
	// 0.2xL+0.4xML+0.4xSR
	limits.byLikes = int64(math.Floor(float64(totalLimit) * 0.2))
	limits.byRating = int64(math.Floor(float64(totalLimit) * 0.4))
	limits.byOracle = int64(math.Floor(float64(totalLimit) * 0.4))
	limits.byRandom = totalLimit - limits.byLikes - limits.byRating - limits.byOracle

	return limits
}
