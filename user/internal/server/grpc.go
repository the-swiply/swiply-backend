package server

import (
	"context"
	"errors"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"
	"github.com/the-swiply/swiply-backend/user/internal/domain"
	"github.com/the-swiply/swiply-backend/user/internal/service"
	"github.com/the-swiply/swiply-backend/user/pkg/api/user"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/mail"
)

type GRPCServer struct {
	user.UnimplementedUserServer
	*grpc.Server

	userService *service.UserService
}

func NewGRPCServer(userService *service.UserService) *GRPCServer {
	srv := &GRPCServer{
		userService: userService,
	}

	opts := []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_recovery.UnaryServerInterceptor(withLogAndRecover()),
		)),
	}
	srv.Server = grpc.NewServer(opts...)
	user.RegisterUserServer(srv.Server, srv)

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

func (g *GRPCServer) SendAuthorizationCode(ctx context.Context, req *user.SendAuthorizationCodeRequest) (*user.SendAuthorizationCodeResponse, error) {
	_, err := mail.ParseAddress(req.GetEmail())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid email")
	}

	err = g.userService.SendAuthorizationCode(ctx, req.GetEmail())
	if errors.Is(err, domain.ErrResendIsNotAllowed) {
		return nil, status.Error(codes.Unavailable, err.Error())
	}
	if err != nil {
		errMsg := fmt.Sprintf("can't send auth code: %v", err.Error())
		loggy.Errorln(errMsg)
		return nil, status.Error(codes.Internal, errMsg)
	}

	return &user.SendAuthorizationCodeResponse{}, err
}

func (g *GRPCServer) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "no incoming meta")
	}

	fingerprint := createFingerprintFromMeta(md)
	tokens, err := g.userService.Login(ctx, req.GetEmail(), req.GetCode(), fingerprint)
	if errors.Is(err, domain.ErrCodeIsIncorrect) || errors.Is(err, domain.ErrTooMuchAttempts) {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	if err != nil {
		errMsg := fmt.Sprintf("can't login user: %v", err.Error())
		loggy.Errorln(errMsg)
		return nil, status.Error(codes.Internal, errMsg)
	}

	return &user.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (g *GRPCServer) Refresh(ctx context.Context, req *user.RefreshRequest) (*user.RefreshResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "no incoming meta")
	}

	fingerprint := createFingerprintFromMeta(md)
	tokens, err := g.userService.RefreshTokens(ctx, req.GetRefreshToken(), fingerprint)
	if errors.Is(err, domain.ErrValidateToken) {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	if err != nil {
		errMsg := fmt.Sprintf("can't refresh tokens: %v", err.Error())
		loggy.Errorln(errMsg)
		return nil, status.Error(codes.Internal, errMsg)
	}

	return &user.RefreshResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
