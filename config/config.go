package config

import (
	"fmt"
	"os"
	"strconv"
)

type (
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"log"`
		PG   `yaml:"postgres"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level" env:"LOG_LEVEL"`
	}

	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"true"                 env:"PG_URL"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	cfg.App.Name = os.Getenv("APP_NAME")
	cfg.App.Version = os.Getenv("APP_VERSION")
	cfg.HTTP.Port = os.Getenv("HTTP_PORT")
	if cfg.HTTP.Port == "" {
		cfg.HTTP.Port = "8080"
	}
	cfg.Log.Level = os.Getenv("LOG_LEVEL")
	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}
	cfg.PG.URL = os.Getenv("PG_URL")
	poolMaxStr := os.Getenv("PG_POOL_MAX")
	if poolMaxStr != "" {
		poolMax, err := strconv.Atoi(poolMaxStr)
		if err == nil {
			cfg.PG.PoolMax = poolMax
		}
	} else {
		cfg.PG.PoolMax = 10
	}
	
	if cfg.PG.URL == "" {
		return nil, fmt.Errorf("PG_URL is required")
	}

	return cfg, nil
}
