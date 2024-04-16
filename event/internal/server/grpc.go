package server

import (
	"context"
	"errors"
	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/the-swiply/swiply-backend/event/internal/converter"
	"github.com/the-swiply/swiply-backend/event/internal/domain"
	"github.com/the-swiply/swiply-backend/event/internal/service"
	"github.com/the-swiply/swiply-backend/event/pkg/api/event"
	"github.com/the-swiply/swiply-backend/pkg/houston/grut"
	"github.com/the-swiply/swiply-backend/pkg/houston/tracy"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	events, err := g.eventService.GetEvents(ctx, req.GetLimit(), req.GetOffset())
	if err != nil {
		return nil, grut.InternalError("can't get user events", err)
	}

	return converter.EventsToGetEvents(events), nil
}

func (g *GRPCServer) GetUserOwnEvents(ctx context.Context, _ *event.GetUserOwnEventsRequest) (*event.GetUserOwnEventsResponse, error) {
	ctx, span := tracy.Start(ctx)
	defer span.End()

	events, err := g.eventService.GetUserOwnEvents(ctx)
	if err != nil {
		return nil, grut.InternalError("can't get user events", err)
	}

	return converter.EventsToGetUserOwnEvents(events), nil
}

func (g *GRPCServer) GetUserMembershipEvents(ctx context.Context, req *event.GetUserMembershipEventsRequest) (*event.GetUserMembershipEventsResponse, error) {
	ctx, span := tracy.Start(ctx)
	defer span.End()

	events, err := g.eventService.GetUserMembershipEvents(ctx)
	if err != nil {
		return nil, grut.InternalError("can't get user events", err)
	}

	return converter.EventsToGetUserMembershipEventsResponse(events), nil
}

func (g *GRPCServer) GetEventMembers(ctx context.Context, req *event.GetEventMembersRequest) (*event.GetEventMembersResponse, error) {
	ctx, span := tracy.Start(ctx)
	defer span.End()

	members, err := g.eventService.GetEventMembers(ctx, req.GetEventId())
	if err != nil {
		return nil, grut.InternalError("can't get user events", err)
	}

	return converter.UserEventStatusToPB(members), nil
}

func (g *GRPCServer) JoinEvent(ctx context.Context, req *event.JoinEventRequest) (*event.JoinEventResponse, error) {
	ctx, span := tracy.Start(ctx)
	defer span.End()

	err := g.eventService.JoinEvent(ctx, req.GetEventId())
	if err != nil {
		return nil, grut.InternalError("can't join event", err)
	}

	return &event.JoinEventResponse{}, nil
}

func (g *GRPCServer) AcceptEventJoin(ctx context.Context, req *event.AcceptEventJoinRequest) (*event.AcceptEventJoinResponse, error) {
	ctx, span := tracy.Start(ctx)
	defer span.End()

	userIDParsed, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id format")
	}

	err = g.eventService.AcceptEventJoin(ctx, req.GetEventId(), userIDParsed)
	if errors.Is(err, domain.ErrEntityIsNotExists) {
		return nil, status.Error(codes.InvalidArgument, "no such event")
	}
	if err != nil {
		return nil, grut.InternalError("can't join event", err)
	}

	return &event.AcceptEventJoinResponse{}, nil
}
