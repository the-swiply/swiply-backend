package glsync

import "time"

type RedisSyncerConfig struct {
	Addr     string
	Password string
	DB       int

	LockExpiration time.Duration
}
