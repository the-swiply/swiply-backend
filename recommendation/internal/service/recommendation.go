package service

import (
	"context"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
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
	SaveUserRecommendationHistory(ctx context.Context, userID uuid.UUID, recommended []uuid.UUID) error
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
}

func (r *RecommendationService) Recommend(ctx context.Context, count int64) ([]uuid.UUID, error) {
	userID := auf.ExtractUserIDFromContext[uuid.UUID](ctx)

	recLimits := calculateRecommendationLimits(count)
	notRandomRescSlice := make([]uuid.UUID, 0, count)

	byLikesRecs, err := r.recRepo.GetGetRecommendationsByLikes(ctx, userID, recLimits.byLikes)
	if err != nil {
		return nil, fmt.Errorf("can't get recommendations by likes: %w", err)
	}
	notRandomRescSlice = append(notRandomRescSlice, byLikesRecs...)

	byRatingRecs, err := r.recRepo.GetRecommendationsByRating(ctx, userID, defaultMaxRatingDelta, r.cfg.FreezeHoursForRecommendation, recLimits.byLikes)
	if err != nil {
		return nil, fmt.Errorf("can't get recommendations by rating: %w", err)
	}
	notRandomRescSlice = append(notRandomRescSlice, byRatingRecs...)

	byOracleRecs, err := r.recRepo.GetRecommendationsByOracle(ctx, userID, r.cfg.FreezeHoursForRecommendation, recLimits.byOracle)
	if err != nil {
		return nil, fmt.Errorf("can't get recommendations by oracle: %w", err)
	}
	notRandomRescSlice = append(notRandomRescSlice, byOracleRecs...)

	recsSet := mapset.NewThreadUnsafeSet[uuid.UUID](notRandomRescSlice...)
	randomPartSize := count - int64(recsSet.Cardinality())

	byRandomRecs, err := r.recRepo.GetRecommendationsByRandom(ctx, userID, r.cfg.FreezeHoursForRecommendation, randomPartSize)
	if err != nil {
		return nil, fmt.Errorf("can't get recommendations by random: %w", err)
	}

	for _, randRec := range byRandomRecs {
		recsSet.Add(randRec)
	}

	recsSlice := recsSet.ToSlice()

	err = r.recRepo.SaveUserRecommendationHistory(ctx, userID, recsSlice)
	if err != nil {
		return nil, fmt.Errorf("can't save recommendation to history: %w", err)
	}

	return recsSlice, nil
}

func calculateRecommendationLimits(totalLimit int64) recommendationLimits {
	limits := recommendationLimits{}

	// Default rule is
	// 0.2xL+0.4xML+0.4xSR
	limits.byLikes = int64(math.Floor(float64(totalLimit) * 0.2))
	limits.byRating = int64(math.Floor(float64(totalLimit) * 0.4))
	limits.byOracle = int64(math.Floor(float64(totalLimit) * 0.4))

	return limits
}
