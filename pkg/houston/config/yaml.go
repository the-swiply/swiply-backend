package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/the-swiply/swiply-backend/pkg/houston/stage"
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

const (
	defaultConfigPath = "configs"
	configPathEnvKey  = "CONFIG_PATH"

	devConfigName  = "values-dev"
	prodConfigName = "values-prod"
)

var (
	ErrUnknownStage = errors.New("unknown stage")
)

func ReadYAML() error {
	switch {
	case stage.IsDev():
		viper.SetConfigName(devConfigName)
	case stage.IsProd():
		viper.SetConfigName(prodConfigName)
	default:
		return ErrUnknownStage
	}

	cfgPath := os.Getenv(configPathEnvKey)
	if cfgPath == "" {
		cfgPath = defaultConfigPath
	}

	viper.SetConfigType("yaml")
	viper.AddConfigPath(cfgPath)
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("can't read in yaml config: %w", err)
	}

	return nil
}

func ParseYAML(out any) error {
	var cfgName string
	switch {
	case stage.IsDev():
		cfgName = devConfigName
	case stage.IsProd():
		cfgName = prodConfigName
	default:
		return ErrUnknownStage
	}

	cfgPath := os.Getenv(configPathEnvKey)
	if cfgPath == "" {
		cfgPath = defaultConfigPath
	}
	data, err := os.ReadFile(fmt.Sprintf("%s.%s", path.Join(cfgPath, cfgName), "yaml"))
	if err != nil {
		return fmt.Errorf("can't read config: %w", err)
	}

	if err = yaml.Unmarshal(data, out); err != nil {
		return fmt.Errorf("can't unmarshal yaml: %w", err)
	}

	return nil
}
