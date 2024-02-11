package service

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
