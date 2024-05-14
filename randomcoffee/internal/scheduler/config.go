package scheduler

type RedisCronConfig struct {
	Addr     string
	Password string
	DB       int

	RandomCoffeeTriggerCron string
}
