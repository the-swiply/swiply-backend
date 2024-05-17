package app

type Config struct {
	App      Application `yaml:"app"`
	GRPC     `yaml:"grpc"`
	HTTP     `yaml:"http"`
	Swagger  `yaml:"swagger"`
	Redis    `yaml:"redis"`
	Postgres `yaml:"postgres"`
}

type Application struct {
	NumOfMessageSenderWorkers      int64 `yaml:"num_of_message_sender_workers"`
	ChatLockExpirationMilliseconds int64 `yaml:"chat_lock_expiration_milliseconds"`
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

type Redis struct {
	Addr          string  `yaml:"addr"`
	DB            RedisDB `yaml:"db"`
	SkipTLSVerify bool    `yaml:"skip_tls_verify"`
	Secure        bool    `yaml:"secure"`
}

type RedisDB struct {
	Sequence       int64
	MessagesPubSub int64
	Syncer         int64
}

type Postgres struct {
	Username string `yaml:"username"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DBName   string `yaml:"db_name"`
	SSLMode  string `yaml:"ssl_mode"`

	MigrationsFolder string `yaml:"migrations_folder"`
}
