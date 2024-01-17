package cache

type RedisDefaultConfig struct {
	Addr     string
	Password string
	DB       int
}

type RedisSequenceConfig struct {
	RedisDefaultConfig
}
