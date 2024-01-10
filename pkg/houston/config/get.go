package config

import (
	"github.com/spf13/viper"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"
	"time"
)

func String(key string) string {
	val := viper.GetString(key)
	if val == "" {
		loggy.Warnln("config value for key", key, "is empty")
	}

	return val
}

func Int(key string) int {
	return viper.GetInt(key)
}

func Seconds(key string) time.Duration {
	return viper.GetDuration(key) * time.Second
}
