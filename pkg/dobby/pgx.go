package dobby

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGXConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
}

func NewPGXPool(ctx context.Context, cfg PGXConfig) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode))

	if err != nil {
		return nil, fmt.Errorf("can't init pgpool: %w", err)
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't ping pg: %w", err)
	}

	return db, nil
}
