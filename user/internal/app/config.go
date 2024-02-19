package app

type Config struct {
	App     Application `yaml:"app"`
	GRPC    `yaml:"grpc"`
	HTTP    `yaml:"http"`
	Swagger `yaml:"swagger"`
	Mailer  `yaml:"mailer"`
	Redis   `yaml:"redis"`
}

type Application struct {
	AuthCodeTTLMinutes                 int64  `yaml:"auth_code_ttl_minutes"`
	AuthCodeSendingMinRetryTimeMinutes int64  `yaml:"auth_code_sending_min_retry_time_minutes"`
	MaxInvalidCodeAttempts             int64  `yaml:"max_invalid_code_attempts"`
	AccessTokenTTLMinutes              int64  `yaml:"access_token_ttl_minutes"`
	RefreshTokenTTLHours               int64  `yaml:"refresh_token_ttl_hours"`
	JobTimeoutSeconds                  int64  `yaml:"job_timeout_seconds"`
	UUIDNamespace                      string `yaml:"uuid_namespace"`
}

type GRPC struct {
	Addr string `yaml:"addr"`
}

type HTTP struct {
	Addr string `yaml:"addr"`
}

type Swagger struct {
	Path string `yaml:"path"`
}

type Mailer struct {
	SMTPAddr    string `yaml:"smtp_addr"`
	SenderEmail string `yaml:"sender_email"`

	SendTimeoutSeconds    int64 `yaml:"send_timeout_seconds"`
	AfterSendPauseSeconds int64 `yaml:"after_send_pause_seconds"`
}

type Redis struct {
	Addr string  `yaml:"addr"`
	DB   RedisDB `yaml:"db"`
}

type RedisDB struct {
	Codes       int64
	Tokens      int64
	MailerQueue int64
}
