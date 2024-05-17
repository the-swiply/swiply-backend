package pubsub

type RedisPubSubConfig struct {
	Addr          string
	Password      string
	DB            int
	SkipTLSVerify bool
	Secure        bool

	ChannelName string
}
