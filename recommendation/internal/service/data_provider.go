package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/the-swiply/swiply-backend/pkg/houston/dobby"
	"github.com/the-swiply/swiply-backend/recommendation/internal/domain"
)

type DataProviderRepository interface {
	GetLastProfileUpdate(ctx context.Context) (time.Time, error)
	GetLastInteractionUpdate(ctx context.Context) (time.Time, error)
	UpdateLastProfileUpdate(ctx context.Context, ts time.Time) error
	UpdateLastInteractionUpdate(ctx context.Context, ts time.Time) error
	UpsertProfiles(ctx context.Context, profiles []domain.Profile) error
	AddInteractions(ctx context.Context, interactions []domain.Interaction) error
	CalculateRatings(ctx context.Context) (map[string]float64, error)
	UpdateStatistics(ctx context.Context, ratings map[string]float64) error
}

type ProfileClient interface {
	GetProfiles(ctx context.Context, after time.Time) ([]domain.Profile, error)
	GetInteractions(ctx context.Context, after time.Time) ([]domain.Interaction, error)
}

type OracleClient interface {
	RetrainLFMv1(ctx context.Context) error
}

type Transactor interface {
	WithinTransaction(ctx context.Context, f func(ctx context.Context) error, opts dobby.TxOptions) error
}

type DataProviderService struct {
	cfg           DataProviderConfig
	dpRepo        DataProviderRepository
	transactor    Transactor
	oracleClient  OracleClient
	profileClient ProfileClient
}

func NewDataProviderService(cfg DataProviderConfig, dpRepo DataProviderRepository, transactor Transactor,
	oracleClient OracleClient, profileClient ProfileClient) *DataProviderService {

	syncPrepareUpdateCh := make(chan struct{})
	close(syncPrepareUpdateCh)

	return &DataProviderService{
		cfg:           cfg,
		dpRepo:        dpRepo,
		transactor:    transactor,
		oracleClient:  oracleClient,
		profileClient: profileClient,
	}
}

func (d *DataProviderService) PrepareRecommendationData(ctx context.Context) error {
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

	interactions, err := d.profileClient.GetInteractions(ctx, lastInteractionUpdate)
	if err != nil {
		return fmt.Errorf("can't get interactions: %w", err)
	}

	err = d.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		err = d.dpRepo.UpsertProfiles(txCtx, profiles)
		if err != nil {
			return fmt.Errorf("can't upsert profiles: %w", err)
		}

		err = d.dpRepo.AddInteractions(txCtx, interactions)
		if err != nil {
			return fmt.Errorf("can't add interactions: %w", err)
		}

		err = d.updateStatistics(txCtx)
		if err != nil {
			return err
		}

		now := time.Now()

		err = d.dpRepo.UpdateLastProfileUpdate(txCtx, now)
		if err != nil {
			return fmt.Errorf("can't update last time update of profiles: %w", err)
		}

		err = d.dpRepo.UpdateLastInteractionUpdate(txCtx, now)
		if err != nil {
			return fmt.Errorf("can't update last time update of interactions: %w", err)
		}

		return nil
	}, dobby.TxOptions{})

	return err
}

func (d *DataProviderService) updateStatistics(ctx context.Context) error {
	ratings, err := d.dpRepo.CalculateRatings(ctx)
	if err != nil {
		return fmt.Errorf("can't calculate ratings: %w", err)
	}

	err = d.dpRepo.UpdateStatistics(ctx, ratings)
	if err != nil {
		return fmt.Errorf("can't update statistics: %w", err)
	}

	return nil
}

func (d *DataProviderService) UpdateOracleData(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	err := d.oracleClient.RetrainLFMv1(ctx)
	if errors.Is(err, context.DeadlineExceeded) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("can't retrain model: %w", err)
	}

	return nil
}
