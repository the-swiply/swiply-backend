package main

import (
	"context"

	"github.com/the-swiply/swiply-backend/pkg/houston/config"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"
	"github.com/the-swiply/swiply-backend/pkg/houston/runner"

	"github.com/the-swiply/swiply-backend/notification/internal/app"
)

const (
	AppNameConfigKey                 = "app.name"
	JaegerAddrConfigKey              = "jaeger.addr"
	GracefulShutdownTimeoutConfigKey = "app.graceful_shutdown_timeout_seconds"
)

func main() {
	ctx := context.Background()
	loggy.InitDefault()

	err := config.ReadYAML()
	if err != nil {
		loggy.Fatal("can't read and parse config:", err)
	}

	cfg := new(app.Config)
	err = config.ParseYAML(&cfg)
	if err != nil {
		loggy.Fatal("can't read and parse config:", err)
	}

	a := app.NewApp(
		cfg,
		runner.NewRunnerV1(
			config.String(AppNameConfigKey),
			config.String(JaegerAddrConfigKey),
		),
	)

	runner.Start(
		ctx,
		a,
		runner.WithGracefulShutdown(config.Seconds(GracefulShutdownTimeoutConfigKey)),
		runner.WithPanicRecovery(),
	)
}
