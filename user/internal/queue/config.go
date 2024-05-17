package queue

import "time"

type MailerConfig struct {
	RedisAddr          string
	RedisPassword      string
	RedisDB            int
	RedisSkipTLSVerify bool
	RedisSecure        bool

	SendTimeout          time.Duration
	AfterSendWorkerPause time.Duration
}
