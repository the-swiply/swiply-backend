package app

type Config struct {
	App      Application `yaml:"app"`
	GRPC     `yaml:"grpc"`
	HTTP     `yaml:"http"`
	Swagger  `yaml:"swagger"`
	Chat     `yaml:"chat"`
	Postgres `yaml:"postgres"`
	S3       `yaml:"s3"`
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

type Postgres struct {
	Username string `yaml:"username"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DBName   string `yaml:"db_name"`
	SSLMode  string `yaml:"ssl_mode"`

	MigrationsFolder string `yaml:"migrations_folder"`
}

type Chat struct {
	Addr string `yaml:"addr"`
}

type S3 struct {
	Addr       string `yaml:"addr"`
	BucketName string `yaml:"bucket_name"`
	AccessKey  string `yaml:"access_key"`
	Secure     bool   `yaml:"secure"`
}
