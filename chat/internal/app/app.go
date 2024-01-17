package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/the-swiply/swiply-backend/chat/internal/cache"
	"github.com/the-swiply/swiply-backend/chat/internal/server"
	"github.com/the-swiply/swiply-backend/chat/internal/service"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"
	"github.com/the-swiply/swiply-backend/pkg/houston/runner"
	"go.uber.org/multierr"
	"net"
	"net/http"
	"os"
)

const (
	sequenceRedisDB = 0
)

type App struct {
	runner.RunStopper
	cfg *Config

	grpcServer *server.GRPCServer
	httpServer *server.HTTPServer

	redisSequenceCache *cache.RedisSequenceCache

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

	rdbSequence, err := cache.NewRedisSequenceCache(ctx, cache.RedisSequenceConfig{
		RedisDefaultConfig: cache.RedisDefaultConfig{
			Addr:     a.cfg.Redis.Addr,
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       sequenceRedisDB,
		},
	})
	if err != nil {
		return fmt.Errorf("can't init redis cache: %w", err)
	}

	a.redisSequenceCache = rdbSequence

	chatSvc := service.NewChatService()

	errCh := make(chan error, 2)

	go func() {
		if err = a.runGRPCServer(chatSvc); err != nil {
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

func (a *App) runGRPCServer(chatSvc *service.ChatService) error {
	srv := server.NewGRPCServer(chatSvc)
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
		a.redisSequenceCache.Stop(ctx),
		a.RunStopper.Stop(ctx),
	)

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
