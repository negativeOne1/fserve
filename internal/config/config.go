package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Log     LogConfig
	HTTP    HTTPConfig
	Storage StorageConfig
	Secret  string `envconfig:"SECRET"`
}

type StorageConfig struct {
	BasePath string `envconfig:"STORAGE_BASE_PATH" default:"/tmp/fserve"`
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

	err := envconfig.Process("FSERVE", &c)
	if err != nil {
		return c, err
	}
	return c, nil
}
