package scheduler

type RedisCronConfig struct {
	Addr          string
	Password      string
	DB            int
	SkipTLSVerify bool
	Secure        bool

	RandomCoffeeTriggerCron string
}
