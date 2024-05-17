package server

import (
	"context"

	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/the-swiply/swiply-backend/pkg/houston/auf"
	"github.com/the-swiply/swiply-backend/pkg/houston/grut"
	"github.com/the-swiply/swiply-backend/pkg/houston/tracy"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/the-swiply/swiply-backend/notification/internal/domain"
	"github.com/the-swiply/swiply-backend/notification/internal/service"
	"github.com/the-swiply/swiply-backend/notification/pkg/api/notification"
)

type GRPCServer struct {
	notification.UnimplementedNotificationServer
	*grpc.Server

	notificationService *service.NotificationService
}

func (g *GRPCServer) Subscribe(ctx context.Context, req *notification.SubscribeRequest) (*notification.SubscribeResponse, error) {
	ctx, span := tracy.Start(ctx)
	defer span.End()

	if err := g.notificationService.Subscribe(ctx, domain.Notification{
		ID:          auf.ExtractUserIDFromContext[uuid.UUID](ctx),
		DeviceToken: req.DeviceToken,
	}); err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &notification.SubscribeResponse{}, nil
}

func (g *GRPCServer) Unsubscribe(ctx context.Context, req *notification.UnsubscribeRequest) (*notification.UnsubscribeResponse, error) {
	ctx, span := tracy.Start(ctx)
	defer span.End()

	if err := g.notificationService.Unsubscribe(ctx, auf.ExtractUserIDFromContext[uuid.UUID](ctx)); err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &notification.UnsubscribeResponse{}, nil
}

func (g *GRPCServer) Send(ctx context.Context, req *notification.SendRequest) (*notification.SendResponse, error) {
	ctx, span := tracy.Start(ctx)
	defer span.End()

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "can't parse id")
	}

	if err := g.notificationService.Send(ctx, id, req.Content); err == domain.ErrEntityIsNotExists {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &notification.SendResponse{}, nil
}

func NewGRPCServer(notificationService *service.NotificationService) *GRPCServer {
	srv := &GRPCServer{
		notificationService: notificationService,
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
	notification.RegisterNotificationServer(srv.Server, srv)

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
