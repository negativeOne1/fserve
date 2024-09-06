package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Log     LogConfig
	HTTP    HTTPConfig
	Storage StorageConfig
}

type StorageConfig struct {
	BasePath string `envconfig:"STORAGE_BASE_PATH" default:"/tmp/"`
}

type HTTPConfig struct {
	Port int `envconfig:"HTTP_PORT" default:"8080"`
}

type LogConfig struct {
	Level  string `envconfig:"LOG_LEVEL"  default:"info"`
	Format string `envconfig:"LOG_FORMAT" default:"json"`
}

func ReadConfig() (Config, error) {
	c := Config{}

	err := envconfig.Process("", &c)
	if err != nil {
		return c, err
	}
	return c, nil
}
