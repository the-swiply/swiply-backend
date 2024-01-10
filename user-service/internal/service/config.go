package service

import "time"

type UserConfig struct {
	MaxAuthCodeTTLForResend time.Duration
	AccessTokenTTL          time.Duration
	RefreshTokenTTL         time.Duration
	TokenSecret             string
}
