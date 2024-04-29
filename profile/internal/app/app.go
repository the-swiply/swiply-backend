package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/the-swiply/swiply-backend/pkg/houston/dobby"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"
	"github.com/the-swiply/swiply-backend/pkg/houston/runner"
	"go.uber.org/multierr"

	"github.com/the-swiply/swiply-backend/profile/internal/repository"
	"github.com/the-swiply/swiply-backend/profile/internal/server"
	"github.com/the-swiply/swiply-backend/profile/internal/service"
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

	db *pgxpool.Pool
	s3 *minio.Client

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

	if os.Getenv("DISABLE_AUTO_MIGRATION") == "" {
		loggy.Infoln("starting migrations")
		err = dobby.AutoMigratePostgres(a.db, a.cfg.Postgres.MigrationsFolder)
		if err != nil {
			return fmt.Errorf("can't apply migrations: %w", err)
		}
		loggy.Infoln("migration done")
	}

	minioClient, err := minio.New(a.cfg.S3.Addr, &minio.Options{
		Creds:  credentials.NewStaticV4(a.cfg.S3.AccessKey, os.Getenv("PHOTO_STORAGE_SECRET_KEY"), ""),
		Secure: true,
	})
	a.s3 = minioClient

	profileRepo := repository.NewProfileRepository(a.db)
	photoContentRepo := repository.NewPhotoContentRepository(a.cfg.S3.BucketName, a.s3)
	photoRepo := repository.NewPhotoRepository(a.db)

	profileSvc := service.NewProfileService(service.ProfileConfig{}, profileRepo)
	photoSvc := service.NewPhotoService(service.PhotoConfig{}, photoContentRepo, photoRepo)

	errCh := make(chan error, 2)

	go func() {
		if err = a.runGRPCServer(profileSvc, photoSvc); err != nil {
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

func (a *App) runGRPCServer(profileSvc *service.ProfileService, photoSvc *service.PhotoService) error {
	srv := server.NewGRPCServer(profileSvc, photoSvc)
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
