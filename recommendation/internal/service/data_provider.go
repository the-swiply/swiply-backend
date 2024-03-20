package service

import (
	"context"
	"fmt"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"
	"github.com/the-swiply/swiply-backend/recommendation/internal/domain"
	"time"
)

type DataProviderRepository interface {
	GetLastProfileUpdate(ctx context.Context) (time.Time, error)
	GetLastInteractionUpdate(ctx context.Context) (time.Time, error)
	UpdateLastProfileUpdate(ctx context.Context, ts time.Time) error
	UpdateLastInteractionUpdate(ctx context.Context, ts time.Time) error
	UpsertProfiles(ctx context.Context, profiles []domain.Profile) error
	UpsertInteractions(ctx context.Context, interactions []domain.Interaction) error
	CalculateRatings(ctx context.Context) (map[string]float64, error)
	UpdateStatistics(ctx context.Context, ratings map[string]float64) error
}
type ProfileClient interface {
	GetInteractions(ctx context.Context, from time.Time) ([]domain.Interaction, error)
	GetProfiles(ctx context.Context, from time.Time) ([]domain.Profile, error)
}

type OracleClient interface {
	RetrainLFMv1(ctx context.Context) error
}

type DataProviderService struct {
	cfg           DataProviderConfig
	dpRepo        DataProviderRepository
	oracleClient  OracleClient
	profileClient ProfileClient
}

func NewDataProviderService(cfg DataProviderConfig, dpRepo DataProviderRepository,
	oracleClient OracleClient, profileClient ProfileClient) *DataProviderService {
	return &DataProviderService{
		cfg:           cfg,
		dpRepo:        dpRepo,
		oracleClient:  oracleClient,
		profileClient: profileClient,
	}
}

func (d *DataProviderService) UpdateStatistic(ctx context.Context) error {
	lastProfileUpdate, err := d.dpRepo.GetLastProfileUpdate(ctx)
	if err != nil {
		return fmt.Errorf("can't get last time update of profiles: %w", err)
	}

	lastInteractionUpdate, err := d.dpRepo.GetLastInteractionUpdate(ctx)
	if err != nil {
		return fmt.Errorf("can't get last time update of interactions: %w", err)
	}

	profiles, err := d.profileClient.GetProfiles(ctx, lastProfileUpdate)
	if err != nil {
		return fmt.Errorf("can't get profiles: %w", err)
	}

	err = d.dpRepo.UpsertProfiles(ctx, profiles)
	if err != nil {
		return fmt.Errorf("can't upsert pofiles")
	}

	interactions, err := d.profileClient.GetInteractions(ctx, lastInteractionUpdate)
	if err != nil {
		return fmt.Errorf("can't get interactions: %w", err)
	}

	err = d.dpRepo.UpsertInteractions(ctx, interactions)
	if err != nil {
		return fmt.Errorf("can't upsert interactions")
	}

	go d.updateStatistics(ctx)

	err = d.dpRepo.UpdateLastProfileUpdate(ctx, time.Now())
	if err != nil {
		return fmt.Errorf("can't udpate last time update of profiles: %w", err)
	}

	err = d.dpRepo.UpdateLastInteractionUpdate(ctx, time.Now())
	if err != nil {
		return fmt.Errorf("can't update last time update of interactions: %w", err)
	}

	return nil
}

func (d *DataProviderService) updateStatistics(ctx context.Context) {
	ratings, err := d.dpRepo.CalculateRatings(ctx)
	if err != nil {
		loggy.Error("can't calculate ratings:", err)
	}

	err = d.dpRepo.UpdateStatistics(ctx, ratings)
	if err != nil {
		loggy.Error("can't update statistics:", err)
	}
}

func (d *DataProviderService) UpdateOracleData(ctx context.Context) error {
	err := d.oracleClient.RetrainLFMv1(ctx)
	if err != nil {
		return fmt.Errorf("can't retrain model: %w", err)
	}

	return nil
}
