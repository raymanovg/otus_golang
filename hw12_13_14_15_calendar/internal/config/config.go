package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var VERSION = "UNKNOWN"

type Config struct {
	Logger LoggerConf
	Server ServerConf
}

type LoggerConf struct {
	Output  []string
	Level   string
	DevMode bool
}

type ServerConf struct {
	Addr string
}

func NewConfig(configFile string) (Config, error) {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}
	cfg := Config{}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return cfg, nil
}
