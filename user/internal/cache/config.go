package cache

import "time"

type RedisDefaultConfig struct {
	Addr          string
	Password      string
	DB            int
	SkipTLSVerify bool
	Secure        bool
}

type RedisCodesConfig struct {
	RedisDefaultConfig
	AuthCodeTTL time.Duration
}

type RedisTokensConfig struct {
	RedisDefaultConfig
	RefreshTokenTTL time.Duration
}
