package tracy

import (
	"context"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"runtime"
)

func Start(ctx context.Context, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	fName := functionCallerName(1)

	return instance.Start(ctx, fName, opts...)
}

func StartWithName(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return instance.Start(ctx, spanName, opts...)
}

func RecordError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}

func functionCallerName(depth int) string {
	fpcs := make([]uintptr, depth)

	n := runtime.Callers(3, fpcs)
	if n == 0 {
		panic("no caller")
	}

	caller := runtime.FuncForPC(fpcs[0] - 1)
	if caller == nil {
		panic("nil caller")
	}

	return caller.Name()
}
