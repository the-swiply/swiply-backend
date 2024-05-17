package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/baobabus/go-apns/apns2"
	"github.com/baobabus/go-apns/cryptox"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/multierr"

	"github.com/the-swiply/swiply-backend/pkg/houston/dobby"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"
	"github.com/the-swiply/swiply-backend/pkg/houston/runner"

	"github.com/the-swiply/swiply-backend/notification/internal/repository"
	"github.com/the-swiply/swiply-backend/notification/internal/server"
	"github.com/the-swiply/swiply-backend/notification/internal/service"
)

const (
	authConfigPath = "configs/authorization.yaml"
)

type App struct {
	runner.RunStopper
	cfg *Config

	grpcServer *server.GRPCServer
	httpServer *server.HTTPServer

	db     *pgxpool.Pool
	apns   *apns2.Client
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

	signingKey, err := cryptox.PKCS8PrivateKeyFromFile(a.cfg.APNS.SigningKeyPath)
	if err != nil {
		return fmt.Errorf("can't read token signing key: %w", err)
	}

	a.apns = &apns2.Client{
		Gateway: apns2.Gateway.Development,
		Signer: &apns2.JWTSigner{
			KeyID:      os.Getenv("KEY_ID"),
			TeamID:     os.Getenv("TEAM_ID"),
			SigningKey: signingKey,
		},
		CommsCfg: apns2.CommsDefault,
		ProcCfg:  apns2.UnlimitedProcConfig,
	}
	err = a.apns.Start(nil)
	if err != nil {
		return fmt.Errorf("can't start apns client: %w", err)
	}

	notificationRepo := repository.NewNotificationRepository(a.db)
	notificationSvc := service.NewNotificationService(service.NotificationServiceConfig{
		Topic: a.cfg.App.Topic,
	}, notificationRepo, a.apns)

	errCh := make(chan error, 2)
	go func() {
		if err = a.runGRPCServer(notificationSvc); err != nil {
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

func (a *App) runGRPCServer(notificationSvc *service.NotificationService) error {
	srv := server.NewGRPCServer(notificationSvc)
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
		a.apns.Stop(),
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
