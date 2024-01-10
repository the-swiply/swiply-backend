package tracy

import (
	"context"
	"fmt"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"time"
)

const (
	tracerName          = "tracy_main"
	initExporterTimeout = 5 * time.Second
	shutdownTimeout     = 5 * time.Second
)

var (
	instance trace.Tracer
)

func Init(ctx context.Context, endpoint string, appName string, opts ...trace.TracerOption) (func(), error) {
	res, err := resource.New(ctx,
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(appName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("can't create resources: %w", err)
	}

	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock()))

	sctx, cancel := context.WithTimeout(ctx, initExporterTimeout)
	defer cancel()

	traceExp, err := otlptrace.New(sctx, traceClient)
	if err != nil {
		return nil, fmt.Errorf("can't init trace exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExp)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	instance = otel.Tracer(tracerName, opts...)
	loggy.Infoln("tracer successfully inited")

	return func() {
		ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()

		if err := traceExp.Shutdown(ctx); err != nil {
			otel.Handle(err)
		}
	}, nil
}
