package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var VERSION = "UNKNOWN"

type (
	Server struct {
		Http HttpServerConf
		Grpc GrpcServerConf
	}
	Config struct {
		Logger LoggerConf
		Server Server
		App    AppConf
	}

	LoggerConf struct {
		Output  []string
		Level   string
		DevMode bool
	}

	HttpServerConf struct {
		Addr string
	}

	GrpcServerConf struct {
		Addr string
	}

	AppConf struct {
		Storage Storage
	}

	Storage struct {
		Name   string
		SQL    SQLStorage
		Memory Memory
	}

	SQLStorage struct {
		DSN          string
		MaxIdleConns int
		MaxOpenConns int
	}

	Memory struct{}
)

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
