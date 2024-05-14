package server

import (
	"context"
	"time"

	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/the-swiply/swiply-backend/pkg/houston/auf"
	"github.com/the-swiply/swiply-backend/pkg/houston/grut"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/the-swiply/swiply-backend/randomcoffee/internal/converter"
	"github.com/the-swiply/swiply-backend/randomcoffee/internal/domain"
	"github.com/the-swiply/swiply-backend/randomcoffee/internal/service"
	"github.com/the-swiply/swiply-backend/randomcoffee/pkg/api/randomcoffee"
)

type GRPCServer struct {
	randomcoffee.UnimplementedRandomCoffeeServer
	*grpc.Server

	meetingService *service.MeetingService
}

func (g *GRPCServer) Create(ctx context.Context, req *randomcoffee.CreateMeetingRequest) (*randomcoffee.CreateMeetingResponse, error) {
	meet := &randomcoffee.Meeting{
		Id:             uuid.New().String(),
		OwnerId:        auf.ExtractUserIDFromContext[uuid.UUID](ctx).String(),
		Start:          req.Start,
		End:            req.End,
		OrganizationId: req.OrganizationId,
		Status:         randomcoffee.MeetingStatus_AWAITING_SCHEDULE,
		CreatedAt:      timestamppb.New(time.Now().UTC()),
	}

	meeting, err := converter.MeetingFromProtoToDomain(meet)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := g.meetingService.Create(ctx, meeting); err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &randomcoffee.CreateMeetingResponse{
		Meeting: meet,
	}, nil
}

func (g *GRPCServer) Delete(ctx context.Context, req *randomcoffee.DeleteMeetingRequest) (*randomcoffee.DeleteMeetingResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "can't parse id")
	}

	if err := g.meetingService.Delete(ctx, id, auf.ExtractUserIDFromContext[uuid.UUID](ctx)); err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &randomcoffee.DeleteMeetingResponse{}, nil
}

func (g *GRPCServer) Update(ctx context.Context, req *randomcoffee.UpdateMeetingRequest) (*randomcoffee.UpdateMeetingResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "can't parse id")
	}

	if err := g.meetingService.Update(ctx, domain.Meeting{
		ID:             id,
		OwnerID:        auf.ExtractUserIDFromContext[uuid.UUID](ctx),
		Start:          req.Start.AsTime(),
		End:            req.Start.AsTime(),
		OrganizationID: req.OrganizationId,
	}); err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &randomcoffee.UpdateMeetingResponse{}, nil
}

func (g *GRPCServer) List(ctx context.Context, _ *randomcoffee.ListMeetingsRequest) (*randomcoffee.ListMeetingsResponse, error) {
	meetings, err := g.meetingService.List(ctx, auf.ExtractUserIDFromContext[uuid.UUID](ctx))
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	resp := &randomcoffee.ListMeetingsResponse{}
	for _, meeting := range meetings {
		resp.Meetings = append(resp.Meetings, converter.MeetingFromDomainToProto(meeting))
	}

	return resp, nil
}

func (g *GRPCServer) Get(ctx context.Context, req *randomcoffee.GetMeetingRequest) (*randomcoffee.GetMeetingResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "can't parse id")
	}

	meeting, err := g.meetingService.Get(ctx, id)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &randomcoffee.GetMeetingResponse{
		Meeting: converter.MeetingFromDomainToProto(meeting),
	}, nil
}

func NewGRPCServer(meetingService *service.MeetingService) *GRPCServer {
	srv := &GRPCServer{
		meetingService: meetingService,
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
	randomcoffee.RegisterRandomCoffeeServer(srv.Server, srv)

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
