package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/the-swiply/swiply-backend/recommendation/internal/domain"
)

func extractUserIDFromContext(ctx context.Context) uuid.UUID {
	return ctx.Value(domain.UserIDKey{}).(uuid.UUID)
}
