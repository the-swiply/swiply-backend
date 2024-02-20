package service

import (
	"context"
	"fmt"
)

type DataProviderRepository interface {
}

type OracleClient interface {
	RetrainLFMv1(ctx context.Context) error
}

type DataProviderService struct {
	cfg          DataProviderConfig
	dpRepo       DataProviderRepository
	oracleClient OracleClient
}

func NewDataProviderService(cfg DataProviderConfig, dpRepo DataProviderRepository, oracleClient OracleClient) *DataProviderService {
	return &DataProviderService{
		cfg:          cfg,
		dpRepo:       dpRepo,
		oracleClient: oracleClient,
	}
}

func (d *DataProviderService) UpdateStatistic(ctx context.Context) error {
	return nil
}

func (d *DataProviderService) UpdateOracleData(ctx context.Context) error {
	err := d.oracleClient.RetrainLFMv1(ctx)
	if err != nil {
		return fmt.Errorf("can't retrain model: %w", err)
	}

	return nil
}
