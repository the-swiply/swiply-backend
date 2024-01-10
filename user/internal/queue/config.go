package queue

import "time"

type MailerConfig struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	SendTimeout          time.Duration
	AfterSendWorkerPause time.Duration
}
