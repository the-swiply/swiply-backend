package service

import (
	"github.com/google/uuid"
	"time"
)

type UserConfig struct {
	MaxAuthCodeTTLForResend time.Duration
	AccessTokenTTL          time.Duration
	RefreshTokenTTL         time.Duration
	TokenSecret             string
	UUIDNamespace           uuid.UUID
}
