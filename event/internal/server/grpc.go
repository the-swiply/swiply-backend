package server

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/the-swiply/swiply-backend/event/internal/converter"
	"github.com/the-swiply/swiply-backend/event/internal/service"
	"github.com/the-swiply/swiply-backend/event/pkg/api/event"
	"github.com/the-swiply/swiply-backend/pkg/houston/grut"
	"github.com/the-swiply/swiply-backend/pkg/houston/tracy"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	event.UnimplementedEventServer
	*grpc.Server

	eventService *service.EventService
}

func NewGRPCServer(eventService *service.EventService) *GRPCServer {
	srv := &GRPCServer{
		eventService: eventService,
	}

	opts := []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_auth.UnaryServerInterceptor(nil),
			grpc_recovery.UnaryServerInterceptor(grut.WithLogAndRecover()),
		)),
	}
	srv.Server = grpc.NewServer(opts...)
	event.RegisterEventServer(srv.Server, srv)

	return srv
}

func (g *GRPCServer) Shutdown(ctx context.Context) error {
	stopCh := make(chan struct{})

	go func() {
		g.Server.GracefulStop()
		close(stopCh)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-stopCh:
		return nil
	}
}

func (g *GRPCServer) CreateEvent(ctx context.Context, req *event.CreateEventRequest) (*event.CreateEventResponse, error) {
	ctx, span := tracy.Start(ctx)
	defer span.End()

	eventID, err := g.eventService.CreateEvent(ctx, converter.EventFromCreateEventRequest(req))
	if err != nil {
		return nil, grut.InternalError("can't create event", err)
	}

	return &event.CreateEventResponse{EventId: eventID}, nil
}

func (g *GRPCServer) UpdateEvent(ctx context.Context, req *event.UpdateEventRequest) (*event.UpdateEventResponse, error) {
	ctx, span := tracy.Start(ctx)
	defer span.End()

	err := g.eventService.UpdateEvent(ctx, converter.EventFromUpdateEventRequest(req))
	if err != nil {
		return nil, grut.InternalError("can't update event", err)
	}

	return &event.UpdateEventResponse{}, nil
}

func (g *GRPCServer) GetEvents(ctx context.Context, req *event.GetEventsRequest) (*event.GetEventsResponse, error) {
	ctx, span := tracy.Start(ctx)
	defer span.End()

	return nil, nil
}

func (g *GRPCServer) GetUserEvents(ctx context.Context, _ *event.GetUserEventsRequest) (*event.GetUserEventsResponse, error) {
	ctx, span := tracy.Start(ctx)
	defer span.End()

	events, err := g.eventService.GetUserEvents(ctx)
	if err != nil {
		return nil, grut.InternalError("can't update event", err)
	}

	return converter.EventsToGetUserEventsResponse(events), nil
}

func (g *GRPCServer) JoinEvent(ctx context.Context, req *event.JoinEventRequest) (*event.JoinEventResponse, error) {
	ctx, span := tracy.Start(ctx)
	defer span.End()

	return nil, nil
}

func (g *GRPCServer) AcceptEventJoin(ctx context.Context, req *event.AcceptEventJoinRequest) (*event.AcceptEventJoinResponse, error) {
	ctx, span := tracy.Start(ctx)
	defer span.End()

	return nil, nil
}
