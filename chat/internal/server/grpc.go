package server

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/the-swiply/swiply-backend/chat/internal/service"
	"github.com/the-swiply/swiply-backend/chat/pkg/api/chat"
	"github.com/the-swiply/swiply-backend/pkg/houston/grut"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	chat.UnimplementedChatServer
	*grpc.Server

	chatService *service.ChatService
}

func NewGRPCServer(chatService *service.ChatService) *GRPCServer {
	srv := &GRPCServer{
		chatService: chatService,
	}

	opts := []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_recovery.UnaryServerInterceptor(grut.WithLogAndRecover()),
		)),
	}
	srv.Server = grpc.NewServer(opts...)
	chat.RegisterChatServer(srv.Server, srv)

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

func (g *GRPCServer) SendMessage(ctx context.Context, req *chat.SendMessageRequest) (*chat.SendMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}

func (g *GRPCServer) GetNextMessages(ctx context.Context, req *chat.GetNextMessagesRequest) (*chat.GetNextMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNextMessages not implemented")
}

func (g *GRPCServer) GetPreviousMessages(ctx context.Context, req *chat.GetPreviousMessagesRequest) (*chat.GetPreviousMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPreviousMessages not implemented")
}
