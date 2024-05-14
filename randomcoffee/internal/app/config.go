package app

import (
	"time"
)

type Config struct {
	App      Application `yaml:"app"`
	GRPC     `yaml:"grpc"`
	HTTP     `yaml:"http"`
	Swagger  `yaml:"swagger"`
	Postgres `yaml:"postgres"`
	Redis    `yaml:"redis"`
}

type Application struct {
	RandomCoffeeTriggerCron string        `yaml:"random_coffee_trigger_cron"`
	MeetingMinInterval      time.Duration `yaml:"meeting_min_interval"`
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
	Addr string `yaml:"addr"`
}
