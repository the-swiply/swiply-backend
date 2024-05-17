package app

type Config struct {
	App      Application `yaml:"app"`
	GRPC     `yaml:"grpc"`
	HTTP     `yaml:"http"`
	Swagger  `yaml:"swagger"`
	Postgres `yaml:"postgres"`
	Redis    `yaml:"redis"`
	Oracle   `yaml:"oracle"`
	Profile  `yaml:"profile"`
}

type Application struct {
	StatisticUpdateCron          string `yaml:"statistic_update_cron"`
	TriggerOracleLearnCron       string `yaml:"trigger_oracle_learn_cron"`
	FreezeHoursForRecommendation int64  `yaml:"freeze_hours_for_recommendation"`
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

type Postgres struct {
	Username string `yaml:"username"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DBName   string `yaml:"db_name"`
	SSLMode  string `yaml:"ssl_mode"`

	MigrationsFolder string `yaml:"migrations_folder"`
}

type Redis struct {
	Addr          string  `yaml:"addr"`
	DB            RedisDB `yaml:"db"`
	SkipTLSVerify bool    `yaml:"skip_tls_verify"`
	Secure        bool    `yaml:"secure"`
}

type RedisDB struct {
	Cron int64
}

type Oracle struct {
	Addr string `yaml:"addr"`
}

type Profile struct {
	Addr string `yaml:"addr"`
}
