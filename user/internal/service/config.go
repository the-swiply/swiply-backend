package service

import (
	"github.com/google/uuid"
	"time"
)

type UserConfig struct {
	MaxAuthCodeTTLForResend time.Duration
	MaxInvalidCodeAttempts  int64
	AccessTokenTTL          time.Duration
	RefreshTokenTTL         time.Duration
	TokenSecret             string
	UUIDNamespace           uuid.UUID
}
