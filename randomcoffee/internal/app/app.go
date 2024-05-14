package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/the-swiply/swiply-backend/pkg/houston/dobby"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"
	"github.com/the-swiply/swiply-backend/pkg/houston/runner"
	"go.uber.org/multierr"

	"github.com/the-swiply/swiply-backend/randomcoffee/internal/algorithm"
	"github.com/the-swiply/swiply-backend/randomcoffee/internal/repository"
	"github.com/the-swiply/swiply-backend/randomcoffee/internal/scheduler"
	"github.com/the-swiply/swiply-backend/randomcoffee/internal/server"
	"github.com/the-swiply/swiply-backend/randomcoffee/internal/service"
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
	err = server.ParseAuthorizationConfig(authConfigPath)
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

	if os.Getenv("DISABLE_AUTO_MIGRATION") == "" {
		loggy.Infoln("starting migrations")
		err = dobby.AutoMigratePostgres(a.db, a.cfg.Postgres.MigrationsFolder)
		if err != nil {
			return fmt.Errorf("can't apply migrations: %w", err)
		}
		loggy.Infoln("migration done")
	}

	meetingRepo := repository.NewMeetingRepository(a.db)
	algo := algorithm.NewRandomCoffeeAlgorithm(algorithm.RandomCoffeeAlgorithmConfig{
		Interval: a.cfg.App.MeetingMinInterval,
	})

	randomCoffeeSvc := service.NewRandomCoffeeService(service.RandomCoffeeConfig{}, algo, meetingRepo)

	rdbCron, err := scheduler.NewRedisCron(scheduler.RedisCronConfig{
		Addr:                    a.cfg.Redis.Addr,
		Password:                os.Getenv("REDIS_PASSWORD"),
		DB:                      cronRedisDB,
		RandomCoffeeTriggerCron: a.cfg.App.RandomCoffeeTriggerCron,
	}, randomCoffeeSvc)
	if err != nil {
		return fmt.Errorf("can't init redis cron scheduler: %w", err)
	}
	a.redisCron = rdbCron

	if err != nil {
		return fmt.Errorf("can't register update statistic task")
	}

	meetingSvc := service.NewMeetingService(service.MeetingConfig{}, meetingRepo)

	errCh := make(chan error, 4)
	go func() {
		err = a.runRedisServer()
		if err != nil && !errors.Is(err, asynq.ErrServerClosed) {
			errCh <- fmt.Errorf("can't run cron server: %w", err)
		} else {
			errCh <- nil
		}
	}()

	go func() {
		err = a.runScheduler()
		if err != nil && !errors.Is(err, asynq.ErrServerClosed) {
			errCh <- fmt.Errorf("can't run cron scheduler: %w", err)
		} else {
			errCh <- nil
		}
	}()

	go func() {
		if err = a.runGRPCServer(meetingSvc); err != nil {
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

func (a *App) runRedisServer() error {
	if err := a.redisCron.RunServer(); err != nil {
		return fmt.Errorf("can't run: %w", err)
	}

	return nil
}

func (a *App) runScheduler() error {
	if err := a.redisCron.RunScheduler(); err != nil {
		return fmt.Errorf("can't run: %w", err)
	}

	return nil
}

func (a *App) runGRPCServer(meetingSvc *service.MeetingService) error {
	srv := server.NewGRPCServer(meetingSvc)
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
