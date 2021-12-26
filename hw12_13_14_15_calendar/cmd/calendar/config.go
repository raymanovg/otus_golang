package main

import (
	"fmt"
	"github.com/spf13/viper"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
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
