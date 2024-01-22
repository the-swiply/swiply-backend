package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/the-swiply/swiply-backend/chat/internal/domain"
	"github.com/the-swiply/swiply-backend/pkg/auf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
)

func jwtAuthFuncGRPC(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token type: expected bearer token")
	}

	claims, err := auf.ValidateJWTAndExtractClaims(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	userID, ok := claims["id"].(string)
	if !ok || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "token has invalid user id")
	}

	userIDParsed, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid id type: uuid expected")
	}

	return context.WithValue(ctx, domain.UserIDKey{}, userIDParsed), nil
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
		claims, err := auf.ValidateJWTAndExtractClaims(token)
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
			rw.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(rw, "invalid id type: uuid expected")
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), domain.UserIDKey{}, userIDParsed))

		next.ServeHTTP(rw, r)
	})
}
