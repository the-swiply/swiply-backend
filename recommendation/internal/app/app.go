package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/the-swiply/swiply-backend/pkg/houston/dobby"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"
	"github.com/the-swiply/swiply-backend/pkg/houston/runner"
	"github.com/the-swiply/swiply-backend/recommendation/internal/repository"
	"github.com/the-swiply/swiply-backend/recommendation/internal/rpclients"
	"github.com/the-swiply/swiply-backend/recommendation/internal/scheduler"
	"github.com/the-swiply/swiply-backend/recommendation/internal/server"
	"github.com/the-swiply/swiply-backend/recommendation/internal/service"
	"go.uber.org/multierr"
	"net"
	"net/http"
	"os"
)

const (
	authConfigPath = "configs/authorization.yaml"

	cronRedisDB = 0
)

type App struct {
	runner.RunStopper
	cfg *Config

	grpcServer *server.GRPCServer
	httpServer *server.HTTPServer

	db        *pgxpool.Pool
	redisCron *scheduler.RedisCron

	stopCh chan struct{}
}

func NewApp(config *Config, runStopperPreset runner.RunStopper) *App {
	return &App{
		RunStopper: runStopperPreset,
		cfg:        config,
		stopCh:     make(chan struct{}),
	}
}

func (a *App) Run(ctx context.Context) error {
	defer close(a.stopCh)

	err := a.RunStopper.Run(ctx)
	if err != nil {
		return err
	}

	server.SetUserJWTSecret(os.Getenv("JWT_SECRET"))
	server.SetS2SJWTSecret(os.Getenv("S2S_JWT_SECRET"))
	err = server.ParseAuthorizationConfig(authConfigPath)
	if err != nil {
		return fmt.Errorf("can't parse auth config: %w", err)
	}
	if err != nil {
		return fmt.Errorf("can't parse auth config: %w", err)
	}

	db, err := dobby.NewPGXPool(ctx, dobby.PGXConfig{
		Username: a.cfg.Postgres.Username,
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     a.cfg.Postgres.Host,
		Port:     a.cfg.Postgres.Port,
		DBName:   a.cfg.Postgres.DBName,
		SSLMode:  a.cfg.Postgres.SSLMode,
	})
	if err != nil {
		return fmt.Errorf("can't init db: %w", err)
	}
	a.db = db

	recRepo := repository.NewRecommendationRepository(a.db)
	dpRepo := repository.NewDataProviderRepository(a.db)

	oracleClient, err := rpclients.NewOracleClient(a.cfg.Oracle.Addr)
	if err != nil {
		return fmt.Errorf("can't get oracle client: %w", err)
	}
	defer oracleClient.CloseConn()

	dpSvc := service.NewDataProviderService(service.DataProviderConfig{}, dpRepo, oracleClient)

	rdbCron, err := scheduler.NewRedisCron(scheduler.RedisCronConfig{
		Addr:                   a.cfg.Redis.Addr,
		Password:               os.Getenv("REDIS_PASSWORD"),
		DB:                     cronRedisDB,
		StatisticUpdateCron:    a.cfg.App.StatisticUpdateCron,
		TriggerOracleLearnCron: a.cfg.App.TriggerOracleLearnCron,
	}, dpSvc)
	if err != nil {
		return fmt.Errorf("can't init redis cron scheduler: %w", err)
	}
	a.redisCron = rdbCron

	if err != nil {
		return fmt.Errorf("can't register update statistic task")
	}

	recSvc := service.NewRecommendationService(service.RecommendationConfig{}, recRepo)

	errCh := make(chan error, 4)
	go func() {
		err = a.runUpdateServer()
		if err != nil && !errors.Is(err, asynq.ErrServerClosed) {
			errCh <- fmt.Errorf("can't run cron server: %w", err)
		} else {
			errCh <- nil
		}
	}()

	go func() {
		err = a.runUpdateScheduler()
		if err != nil && !errors.Is(err, asynq.ErrServerClosed) {
			errCh <- fmt.Errorf("can't run cron scheduler: %w", err)
		} else {
			errCh <- nil
		}
	}()

	go func() {
		if err = a.runGRPCServer(recSvc); err != nil {
			errCh <- fmt.Errorf("can't run grpc server: %w", err)
		} else {
			errCh <- nil
		}
	}()

	go func() {
		err = a.runHTTPServer(ctx)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("can't run http server: %w", err)
		} else {
			errCh <- nil
		}
	}()

	for i := 0; i < cap(errCh); i++ {
		if err = <-errCh; err != nil {
			return err
		}
	}

	return nil
}

func (a *App) runUpdateServer() error {
	if err := a.redisCron.RunServer(); err != nil {
		return fmt.Errorf("can't run: %w", err)
	}

	return nil
}

func (a *App) runUpdateScheduler() error {
	if err := a.redisCron.RunScheduler(); err != nil {
		return fmt.Errorf("can't run: %w", err)
	}

	return nil
}

func (a *App) runGRPCServer(recSvc *service.RecommendationService) error {
	srv := server.NewGRPCServer(recSvc)
	a.grpcServer = srv

	listener, err := net.Listen("tcp", a.cfg.GRPC.Addr)
	if err != nil {
		return fmt.Errorf("can't listen: %w", err)
	}
	defer listener.Close()

	loggy.Infoln("starting grpc server on", a.cfg.GRPC.Addr)

	if err = a.grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("can't serve: %w", err)
	}

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	srv, err := server.NewHTTPServer(ctx, server.HTTPConfig{
		ServeAddr:    a.cfg.HTTP.Addr,
		GRPCEndpoint: a.cfg.GRPC.Addr,
		SwaggerPath:  a.cfg.Swagger.Path,
	})
	if err != nil {
		return err
	}

	a.httpServer = srv

	loggy.Infoln("starting http server on", a.cfg.HTTP.Addr)
	if err = a.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("can't serve: %w", err)
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	err := multierr.Combine(
		a.grpcServer.Shutdown(ctx),
		a.httpServer.Shutdown(ctx),

		a.RunStopper.Stop(ctx),
	)
	a.db.Close()

	if err != nil {
		return err
	}

	select {
	case <-a.stopCh:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
