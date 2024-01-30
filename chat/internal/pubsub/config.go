package pubsub

type RedisPubSubConfig struct {
	Addr        string
	Password    string
	DB          int
	ChannelName string
}
