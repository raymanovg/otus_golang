package main

import "github.com/spf13/viper"

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf
	// TODO
}

type LoggerConf struct {
	Level string
	// TODO
}

func NewConfig(configFile string) (Config, error) {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	return Config{
		Logger: LoggerConf{
			Level: viper.GetString("logger.level"),
		},
	}, nil
}

// TODO
