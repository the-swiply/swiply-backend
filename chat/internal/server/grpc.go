package server

import (
	"context"
	"errors"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/the-swiply/swiply-backend/chat/internal/converter"
	"github.com/the-swiply/swiply-backend/chat/internal/domain"
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
			grpc_auth.UnaryServerInterceptor(jwtAuthFuncGRPC),
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
	err := g.chatService.ReceiveChatMessage(ctx, req.GetChatId(), req.GetContent())
	switch {
	case errors.Is(err, domain.ErrEntityIsNotExists):
		return nil, status.Error(codes.PermissionDenied, "no such chat")
	case errors.Is(err, domain.ErrUserNotInChat):
		return nil, status.Error(codes.PermissionDenied, domain.ErrUserNotInChat.Error())
	case err != nil:
		return nil, grut.InternalError("can't receive message", err)
	}

	return &chat.SendMessageResponse{}, nil
}

func (g *GRPCServer) GetNextMessages(ctx context.Context, req *chat.GetNextMessagesRequest) (*chat.GetNextMessagesResponse, error) {
	if req.GetLimit() < 0 {
		return nil, status.Error(codes.InvalidArgument, "limit must be not negative")
	}

	messages, err := g.chatService.GetNextMessages(ctx, req.GetChatId(), req.GetStartingFrom(), req.GetLimit())
	switch {
	case errors.Is(err, domain.ErrEntityIsNotExists):
		return nil, status.Error(codes.PermissionDenied, "no such chat")
	case errors.Is(err, domain.ErrUserNotInChat):
		return nil, status.Error(codes.PermissionDenied, domain.ErrUserNotInChat.Error())
	case err != nil:
		return nil, grut.InternalError("can't get next messages", err)
	}

	return &chat.GetNextMessagesResponse{
		Messages: converter.MessagesToPB(messages),
	}, nil
}

func (g *GRPCServer) GetPreviousMessages(ctx context.Context, req *chat.GetPreviousMessagesRequest) (*chat.GetPreviousMessagesResponse, error) {
	if req.GetLimit() < 0 {
		return nil, status.Error(codes.InvalidArgument, "limit must be not negative")
	}

	messages, err := g.chatService.GetPreviousMessages(ctx, req.GetChatId(), req.GetStartingFrom(), req.GetLimit())
	switch {
	case errors.Is(err, domain.ErrEntityIsNotExists):
		return nil, status.Error(codes.PermissionDenied, "no such chat")
	case errors.Is(err, domain.ErrUserNotInChat):
		return nil, status.Error(codes.PermissionDenied, domain.ErrUserNotInChat.Error())
	case err != nil:
		return nil, grut.InternalError("can't get previous messages", err)
	}

	return &chat.GetPreviousMessagesResponse{
		Messages: converter.MessagesToPB(messages),
	}, nil
}

func (g *GRPCServer) GetChats(ctx context.Context, req *chat.GetChatsRequest) (*chat.GetChatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChats not implemented")
}
func (g *GRPCServer) GetChatMembers(ctx context.Context, req *chat.GetChatMembersRequest) (*chat.GetChatMembersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChatMembers not implemented")
}
func (g *GRPCServer) LeaveChat(ctx context.Context, req *chat.LeaveChatRequest) (*chat.LeaveChatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LeaveChat not implemented")
}
