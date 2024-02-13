package repository

import "github.com/jackc/pgx/v5/pgxpool"

type DataProviderRepository struct {
	db *pgxpool.Pool
}

func NewDataProviderRepository(db *pgxpool.Pool) *DataProviderRepository {
	return &DataProviderRepository{
		db: db,
	}
}
