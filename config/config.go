package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"logger"`
		PG   `yaml:"postgres"`
		JWT  `yaml:"jwt"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"true" yaml:"url"      env:"PG_URL"`
	}

	JWT struct {
		AccessTokenEXP int    `env-required:"true" yaml:"exp"      env:"JWT_EXP"`
		JWTSecretKey   string `env-required:"true" yaml:"secret"   env:"JWT_SECRET_KEY"`
		Audience       string `env-required:"true" yaml:"audience" env:"JWT_AUDIENCE"`
	}
)

var cfg *Config

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadEnv(cfg)

	return cfg, err
}
