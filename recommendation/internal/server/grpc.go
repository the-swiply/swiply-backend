package server

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/the-swiply/swiply-backend/pkg/houston/grut"
	"github.com/the-swiply/swiply-backend/recommendation/internal/service"
	"github.com/the-swiply/swiply-backend/recommendation/pkg/api/recommendation"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	recommendation.UnimplementedRecommendationServer
	*grpc.Server

	recService *service.RecommendationService
}

func NewGRPCServer(recService *service.RecommendationService) *GRPCServer {
	srv := &GRPCServer{
		recService: recService,
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
	recommendation.RegisterRecommendationServer(srv.Server, srv)

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
