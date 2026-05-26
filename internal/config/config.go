package config

import "github.com/caarlos0/env/v11"

type Config struct {
	AppEnv string `env:"APP_ENV" envDefault:"development"`

	HTTPPort string `env:"HTTP_PORT" envDefault:"8080"`

	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME"`

	JWTSecret string `env:"JWT_SECRET"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
