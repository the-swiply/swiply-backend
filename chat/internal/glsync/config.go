package glsync

import "time"

type RedisSyncerConfig struct {
	Addr          string
	Password      string
	DB            int
	SkipTLSVerify bool
	Secure        bool

	LockExpiration time.Duration
}
