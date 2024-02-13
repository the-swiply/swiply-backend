package service

import "context"

type DataProviderRepository interface {
}

type DataProviderService struct {
	cfg    DataProviderConfig
	dpRepo DataProviderRepository
}

func NewDataProviderService(cfg DataProviderConfig, dpRepo DataProviderRepository) *DataProviderService {
	return &DataProviderService{
		cfg:    cfg,
		dpRepo: dpRepo,
	}
}

func (d *DataProviderService) UpdateStatistic(ctx context.Context) error {
	return nil
}

func (d *DataProviderService) UpdateOracleData(ctx context.Context) error {
	return nil
}
