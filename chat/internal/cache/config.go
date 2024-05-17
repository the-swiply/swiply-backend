package cache

type RedisSequenceConfig struct {
	Addr          string
	Password      string
	DB            int
	SkipTLSVerify bool
	Secure        bool
}
