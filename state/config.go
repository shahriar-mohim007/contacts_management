package state

import (
	"github.com/caarlos0/env/v9"
)

type Config struct {
	ApplicationPort int    `env:"APPLICATION_PORT"`
	DatabaseUrl     string `env:"DATABASE_URL"`
	LogLevel        string `env:"LOG_LEVEL" envDefault:"debug"`
	SecretKey       string `env:"SECRET_KEY"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := env.ParseWithOptions(cfg, env.Options{RequiredIfNoDef: true})
	return cfg, err
}
