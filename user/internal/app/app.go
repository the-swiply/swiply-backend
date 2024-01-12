package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"
	"github.com/the-swiply/swiply-backend/pkg/houston/runner"
	"github.com/the-swiply/swiply-backend/user/internal/cache"
	"github.com/the-swiply/swiply-backend/user/internal/mailer"
	"github.com/the-swiply/swiply-backend/user/internal/queue"
	"github.com/the-swiply/swiply-backend/user/internal/server"
	"github.com/the-swiply/swiply-backend/user/internal/service"
	"go.uber.org/multierr"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	codesRedisDB       = 0
	tokensRedisDB      = 1
	mailerQueueRedisDB = 2
)

type App struct {
	runner.RunStopper
	cfg *Config

	grpcServer *server.GRPCServer
	httpServer *server.HTTPServer

	redisCodeCache  *cache.RedisCodeCache
	redisTokenCache *cache.RedisTokenCache
	mailerQueue     *queue.MailerQueue

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

	rdbCodes, err := cache.NewRedisCodeCache(ctx, cache.RedisCodesConfig{
		RedisDefaultConfig: cache.RedisDefaultConfig{
			Addr:     a.cfg.Redis.Addr,
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       codesRedisDB,
		},
		AuthCodeTTL: time.Minute * time.Duration(a.cfg.App.AuthCodeTTLMinutes),
	})
	a.redisCodeCache = rdbCodes

	rdbTokens, err := cache.NewRedisTokenCache(ctx, cache.RedisTokensConfig{
		RedisDefaultConfig: cache.RedisDefaultConfig{
			Addr:     a.cfg.Redis.Addr,
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       tokensRedisDB,
		},
		RefreshTokenTTL: time.Duration(a.cfg.App.RefreshTokenTTLHours) * time.Hour,
	})
	a.redisTokenCache = rdbTokens

	if err != nil {
		return fmt.Errorf("can't init redis cache: %w", err)
	}

	mailSender, err := mailer.NewSMTPClient(mailer.SMTPConfig{
		SenderEmail:    a.cfg.Mailer.SenderEmail,
		SenderPassword: os.Getenv("SENDER_PASSWORD"),
		Addr:           a.cfg.Mailer.SMTPAddr,
	})
	if err != nil {
		return fmt.Errorf("can't init smtp client: %w", err)
	}

	senderSvc := service.NewSenderService(mailSender)

	mailerQueue := queue.NewMailerQueue(queue.MailerConfig{
		RedisAddr:            a.cfg.Redis.Addr,
		RedisPassword:        os.Getenv("REDIS_PASSWORD"),
		RedisDB:              mailerQueueRedisDB,
		SendTimeout:          time.Duration(a.cfg.Mailer.SendTimeoutSeconds) * time.Second,
		AfterSendWorkerPause: time.Duration(a.cfg.Mailer.AfterSendPauseSeconds) * time.Second,
	}, senderSvc)

	a.mailerQueue = mailerQueue

	userSvc := service.NewUserService(service.UserConfig{
		MaxAuthCodeTTLForResend: a.calculateMaxAuthCodeTTLForResend(),
		MaxInvalidCodeAttempts:  a.cfg.App.MaxInvalidCodeAttempts,
		AccessTokenTTL:          time.Duration(a.cfg.App.AccessTokenTTLMinutes) * time.Minute,
		RefreshTokenTTL:         time.Duration(a.cfg.App.RefreshTokenTTLHours) * time.Hour,
		TokenSecret:             os.Getenv("JWT_SECRET"),
		UUIDNamespace:           uuid.MustParse(a.cfg.App.UUIDNamespace),
	}, a.redisCodeCache, a.redisTokenCache, a.mailerQueue)

	errCh := make(chan error, 3)

	go func() {
		err = a.runMailerQueueServer()
		if err != nil && !errors.Is(err, asynq.ErrServerClosed) {
			errCh <- fmt.Errorf("can't run mailer queue server: %w", err)
		} else {
			errCh <- nil
		}
	}()

	go func() {
		if err = a.runGRPCServer(userSvc); err != nil {
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

func (a *App) runGRPCServer(userSvc *service.UserService) error {
	srv := server.NewGRPCServer(userSvc)
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

func (a *App) runMailerQueueServer() error {
	if err := a.mailerQueue.Run(); err != nil {
		return fmt.Errorf("can't run: %w", err)
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	err := multierr.Combine(
		a.grpcServer.Shutdown(ctx),
		a.httpServer.Shutdown(ctx),
		a.mailerQueue.Stop(ctx),
		a.redisCodeCache.Stop(ctx),
		a.redisTokenCache.Stop(ctx),
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
