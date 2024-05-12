package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/the-swiply/swiply-backend/pkg/houston/auf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/yaml.v3"
)

const (
	authorizationHeader    = "authorization"
	s2sAuthorizationHeader = "s2s-authorization"
)

var (
	s2sJWTSecret  []byte
	userJWTSecret []byte

	authCfg = map[string]HandlerAuthorizationConfig{}
)

type HandlerAuthorizationConfig struct {
	User       bool
	S2S        bool
	Authorized map[string]struct{}
}

func SetS2SJWTSecret(secret string) {
	s2sJWTSecret = []byte(secret)
}

func SetUserJWTSecret(secret string) {
	userJWTSecret = []byte(secret)
}

func ParseAuthorizationConfig(cfgPath string) error {
	type MethodAuthConfig struct {
		User       bool     `yaml:"user"`
		S2S        bool     `yaml:"s2s"`
		Authorized []string `yaml:"authorized"`
	}

	var yamlAuthCfg map[string]MethodAuthConfig

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return fmt.Errorf("can't read config: %w", err)
	}

	if err = yaml.Unmarshal(data, &yamlAuthCfg); err != nil {
		return fmt.Errorf("can't unmarshal yaml: %w", err)
	}

	authCfg = make(map[string]HandlerAuthorizationConfig, len(yamlAuthCfg))
	for handlerName, cfg := range yamlAuthCfg {
		authCfg[handlerName] = HandlerAuthorizationConfig{
			User:       cfg.User,
			S2S:        cfg.S2S,
			Authorized: make(map[string]struct{}, len(cfg.Authorized)),
		}

		for _, allowedService := range cfg.Authorized {
			authCfg[handlerName].Authorized[allowedService] = struct{}{}
		}
	}

	return nil
}

func authFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	var err error
	if handlerAuthCfg, hasAuthCfg := authCfg[fullMethodName]; hasAuthCfg {
		if handlerAuthCfg.S2S {
			ctx, err = s2sAuth(ctx, handlerAuthCfg.Authorized)
			if err != nil {
				return nil, err
			}
		}

		if handlerAuthCfg.User {
			ctx, err = userAuth(ctx)
			if err != nil {
				return nil, err
			}
		}
	}

	return ctx, nil
}

func (p *GRPCServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return authFuncOverride(ctx, fullMethodName)
}

func userAuth(ctx context.Context) (context.Context, error) {
	token, err := authFromMD(ctx, authorizationHeader, "bearer")
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token type: expected bearer token")
	}

	claims, err := auf.ValidateJWTAndExtractClaimsWithCustomSecret(token, userJWTSecret)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	userID, ok := claims["id"].(string)
	if !ok || userID == "" {
		return nil, status.Error(codes.PermissionDenied, "token has invalid user id")
	}

	userIDParsed, err := uuid.Parse(userID)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "invalid id type: uuid expected")
	}

	return auf.AddUserIDToContext(ctx, userIDParsed), nil
}

func s2sAuth(ctx context.Context, authorizedServices map[string]struct{}) (context.Context, error) {
	token, err := authFromMD(ctx, s2sAuthorizationHeader, "bearer")
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token type: expected bearer token")
	}

	claims, err := auf.ValidateJWTAndExtractClaimsWithCustomSecret(token, s2sJWTSecret)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	serviceName, ok := claims["service-name"].(string)
	if !ok || serviceName == "" {
		return nil, status.Error(codes.PermissionDenied, "token has invalid service name")
	}

	if _, serviceIsAuthorized := authorizedServices[serviceName]; !serviceIsAuthorized {
		return nil, status.Error(codes.PermissionDenied, "service has no permissions to access this method")
	}

	return ctx, nil
}

func authMiddlewareHTTP(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(authorizationHeader)
		if authHeader == "" {
			rw.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(rw, "authorization header is not set")
			return
		}

		splits := strings.SplitN(authHeader, " ", 2)
		if len(splits) < 2 || !strings.EqualFold(splits[0], "bearer") {
			rw.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(rw, "invalid token type: expected bearer token")
			return
		}

		token := splits[1]
		claims, err := auf.ValidateJWTAndExtractClaimsWithCustomSecret(token, userJWTSecret)
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(rw, "invallid auth token: ", err)
			return
		}

		userID, ok := claims["id"].(string)
		if !ok || userID == "" {
			rw.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(rw, "token has invalid user id")
			return
		}

		userIDParsed, err := uuid.Parse(userID)
		if err != nil {
			rw.WriteHeader(http.StatusForbidden)
			fmt.Fprint(rw, "invalid id type: uuid expected")
			return
		}

		r = r.WithContext(auf.AddUserIDToContext(r.Context(), userIDParsed))

		next.ServeHTTP(rw, r)
	})
}

func authFromMD(ctx context.Context, header string, expectedScheme string) (string, error) {
	val := metautils.ExtractIncoming(ctx).Get(header)
	if val == "" {
		return "", status.Errorf(codes.Unauthenticated, "Request unauthenticated with "+expectedScheme)
	}
	splits := strings.SplitN(val, " ", 2)
	if len(splits) < 2 {
		return "", status.Errorf(codes.Unauthenticated, "Bad authorization string")
	}
	if !strings.EqualFold(splits[0], expectedScheme) {
		return "", status.Errorf(codes.Unauthenticated, "Request unauthenticated with "+expectedScheme)
	}
	return splits[1], nil
}
