package app

type Config struct {
	App      Application `yaml:"app"`
	GRPC     `yaml:"grpc"`
	HTTP     `yaml:"http"`
	Swagger  `yaml:"swagger"`
	Postgres `yaml:"postgres"`
	APNS     `yaml:"apns"`
}

type Application struct {
	Topic string `yaml:"topic"`
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

type APNS struct {
	SigningKeyPath string `yaml:"signing_key_path"`
}
