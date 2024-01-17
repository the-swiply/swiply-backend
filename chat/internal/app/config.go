package app

type Config struct {
	App     Application `yaml:"app"`
	GRPC    `yaml:"grpc"`
	HTTP    `yaml:"http"`
	Swagger `yaml:"swagger"`
	Redis   `yaml:"redis"`
}

type Application struct {
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
	Addr string `yaml:"addr"`
}
