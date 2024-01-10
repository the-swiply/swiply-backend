package runner

import (
	"context"
	"fmt"
	"github.com/the-swiply/swiply-backend/pkg/houston/prometrics"
	"github.com/the-swiply/swiply-backend/pkg/houston/tracy"
)

type RunnerV1Option func(runner *RunnerV1)

type RunnerV1 struct {
	appName     string
	jaegerAddr  string
	metricsPath string
	metricsAddr string

	shutdownTracer func()
}

func WithMetricsPath(path string) RunnerV1Option {
	return func(runner *RunnerV1) {
		runner.metricsPath = path
	}
}

func WithMetricsAddr(addr string) RunnerV1Option {
	return func(runner *RunnerV1) {
		runner.metricsAddr = addr
	}
}

func NewRunnerV1(appName, jaegerAddr string, opts ...RunnerV1Option) *RunnerV1 {
	runner := &RunnerV1{
		appName:     appName,
		jaegerAddr:  jaegerAddr,
		metricsPath: "/metrics",
		metricsAddr: ":9000",
	}

	for _, opt := range opts {
		opt(runner)
	}

	return runner
}

func (r *RunnerV1) Run(ctx context.Context) error {
	var err error
	r.shutdownTracer, err = tracy.Init(ctx, r.jaegerAddr, r.appName)
	if err != nil {
		return fmt.Errorf("can't init tracer: %w", err)
	}

	prometrics.RunMetricsHandler(r.metricsPath, r.metricsAddr)

	return nil
}

func (r *RunnerV1) Stop(_ context.Context) error {
	if r.shutdownTracer != nil {
		r.shutdownTracer()
	}

	return nil
}
