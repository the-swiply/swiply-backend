package runner

import (
	"context"
	"errors"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"
	"os/signal"
	"syscall"
	"time"
)

type RunStopper interface {
	Run(ctx context.Context) error
	Stop(ctx context.Context) error
}

type StartOption func(c *startConfig)

func WithGracefulShutdown(timeout time.Duration) StartOption {
	return func(c *startConfig) {
		c.gracefulShutdownTimeout = timeout
	}
}

func WithPanicRecovery() StartOption {
	return func(c *startConfig) {
		c.recoverPanic = true
	}
}

func Start(ctx context.Context, app RunStopper, opts ...StartOption) {
	cfg := newStartConfig(opts...)

	stopCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if cfg.gracefulShutdownTimeout > 0 {
		defer exit(ctx, app, cfg.gracefulShutdownTimeout)
	}

	go func() {
		if cfg.recoverPanic {
			defer func() {
				if v := recover(); v != nil {
					loggy.Errorln("panic in app.Run:", v)
					stop()
				}
			}()
		}

		err := app.Run(ctx)
		if err != nil {
			loggy.Fatal("can't run app:", err)
		}
	}()

	<-stopCtx.Done()
}

func exit(ctx context.Context, app RunStopper, gracefulShutdownTimeout time.Duration) {
	ctx, cancel := context.WithTimeout(ctx, gracefulShutdownTimeout)
	defer cancel()

	loggy.Infoln("stopping app gracefully")

	err := app.Stop(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			loggy.Warnln("graceful stop time is out. hard-shutting down app")
		} else {
			loggy.Errorln("error during graceful shutdown:", err)
		}
	} else {
		loggy.Infoln("app gracefully stopped")
	}
}

type startConfig struct {
	gracefulShutdownTimeout time.Duration
	recoverPanic            bool
}

func newStartConfig(opts ...StartOption) *startConfig {
	cfg := &startConfig{
		gracefulShutdownTimeout: 0,
		recoverPanic:            false,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}
