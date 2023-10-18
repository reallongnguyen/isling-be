package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App     `yaml:"app"`
		HTTP    `yaml:"http"`
		Log     `yaml:"logger"`
		PG      `yaml:"postgres"`
		JWT     `yaml:"jwt"`
		GORSE   `yaml:"gorse"`
		Surreal `yaml:"surreal"`
		Redis   `yaml:"redis"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
		ENV     string `env-required:"true" yaml:"env" env:"APP_ENV"`
	}

	// HTTP -.
	HTTP struct {
		Port                        string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		RateLimit                   int    `env-required:"true" yaml:"rate_limit" env:"HTTP_RATE_LIMIT"`
		RateLimitUserActivitiesPost int    `env-required:"true" yaml:"rate_limit_user_activities_post" env:"HTTP_RATE_LIMIT_USER_ACTIVITIES_POST"`
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

	Surreal struct {
		URL      string `env-required:"true" yaml:"url"     env:"SURREAL_URL"`
		NS       string `env-required:"true" yaml:"ns"      env:"SURREAL_NS"`
		DB       string `env-required:"true" yaml:"db"      env:"SURREAL_DB"`
		User     string `env-required:"true" env:"SURREAL_USER"`
		Password string `env-required:"true" env:"SURREAL_PASS"`
	}

	JWT struct {
		AccessTokenEXP int    `env-required:"true" yaml:"exp"      env:"JWT_EXP"`
		JWTSecretKey   string `env-required:"true" yaml:"secret"   env:"JWT_SECRET_KEY"`
		Audience       string `env-required:"true" yaml:"audience" env:"JWT_AUDIENCE"`
	}

	GORSE struct {
		URL    string `env-required:"true" yaml:"url" env:"GORSE_URL"`
		APIKey string `env:"GORSE_SERVER_API_KEY"`
	}

	Redis struct {
		URL string `env-required:"true" yaml:"url" env:"REDIS_URL"`
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
