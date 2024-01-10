package app

type Config struct {
	App  Application `yaml:"app"`
	HTTP `yaml:"http"`
}

type Application struct {
}

type HTTP struct {
	Addr string `yaml:"addr"`
}
