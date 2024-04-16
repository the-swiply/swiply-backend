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
	"github.com/the-swiply/swiply-backend/pkg/houston/tracy"
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
			grpc_auth.UnaryServerInterceptor(nil),
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
	ctx, span := tracy.Start(ctx)
	defer span.End()

	err := g.chatService.ReceiveChatMessage(ctx, req.GetChatId(), req.GetContent())
	switch {
	case errors.Is(err, domain.ErrEntityIsNotExists):
		return nil, status.Error(codes.PermissionDenied, "no such chat")
	case errors.Is(err, domain.ErrUserNotInChat):
		return nil, status.Error(codes.InvalidArgument, domain.ErrUserNotInChat.Error())
	case err != nil:
		return nil, grut.InternalError("can't receive message", err)
	}

	return &chat.SendMessageResponse{}, nil
}

func (g *GRPCServer) GetNextMessages(ctx context.Context, req *chat.GetNextMessagesRequest) (*chat.GetNextMessagesResponse, error) {
	ctx, span := tracy.Start(ctx)
	defer span.End()

	if req.GetLimit() < 0 {
		return nil, status.Error(codes.InvalidArgument, "limit must be not negative")
	}

	messages, err := g.chatService.GetNextMessages(ctx, req.GetChatId(), req.GetStartingFrom(), req.GetLimit())
	switch {
	case errors.Is(err, domain.ErrEntityIsNotExists):
		return nil, status.Error(codes.PermissionDenied, "no such chat")
	case errors.Is(err, domain.ErrUserNotInChat):
		return nil, status.Error(codes.InvalidArgument, domain.ErrUserNotInChat.Error())
	case err != nil:
		return nil, grut.InternalError("can't get next messages", err)
	}

	return &chat.GetNextMessagesResponse{
		Messages: converter.MessagesToPB(messages),
	}, nil
}

func (g *GRPCServer) GetPreviousMessages(ctx context.Context, req *chat.GetPreviousMessagesRequest) (*chat.GetPreviousMessagesResponse, error) {
	ctx, span := tracy.Start(ctx)
	defer span.End()

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

func (g *GRPCServer) GetChats(ctx context.Context, _ *chat.GetChatsRequest) (*chat.GetChatsResponse, error) {
	ctx, span := tracy.Start(ctx)
	defer span.End()

	chats, err := g.chatService.GetUserChats(ctx)
	if err != nil {
		return nil, grut.InternalError("can't get user's chats", err)
	}

	return &chat.GetChatsResponse{
		Chats: converter.ChatsToPB(chats),
	}, nil
}

func (g *GRPCServer) LeaveChat(ctx context.Context, req *chat.LeaveChatRequest) (*chat.LeaveChatResponse, error) {
	ctx, span := tracy.Start(ctx)
	defer span.End()

	err := g.chatService.LeaveChat(ctx, req.GetChatId())
	switch {
	case errors.Is(err, domain.ErrEntityIsNotExists):
		return nil, status.Error(codes.PermissionDenied, "no such chat")
	case errors.Is(err, domain.ErrUserNotInChat):
		return nil, status.Error(codes.InvalidArgument, domain.ErrUserNotInChat.Error())
	case err != nil:
		return nil, grut.InternalError("user can't leave chat", err)
	}

	return &chat.LeaveChatResponse{}, nil
}

func (g *GRPCServer) CreateChat(ctx context.Context, req *chat.CreateChatRequest) (*chat.CreateChatResponse, error) {
	ctx, span := tracy.Start(ctx)
	defer span.End()

	members, err := converter.StringsToUUIDs(req.GetMembers())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	chatID, err := g.chatService.CreateChat(ctx, members)
	if err != nil {
		return nil, grut.InternalError("can't create chat", err)
	}

	return &chat.CreateChatResponse{ChatId: chatID}, nil
}

func (g *GRPCServer) AddChatMembers(ctx context.Context, req *chat.AddChatMembersRequest) (*chat.AddChatMembersResponse, error) {
	ctx, span := tracy.Start(ctx)
	defer span.End()

	newMembers, err := converter.StringsToUUIDs(req.GetMembers())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = g.chatService.AddChatMembers(ctx, req.GetChatId(), newMembers)
	if err != nil {
		return nil, grut.InternalError("can't add members", err)
	}

	return &chat.AddChatMembersResponse{}, nil
}
