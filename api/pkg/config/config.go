package config

import "github.com/caarlos0/env"

type Config struct {
	Database Database
	CORS     Cors
}

func LoadConfig() (*Config, error) {
	var cfg Config

	var database Database
	if err := env.Parse(&database); err != nil {
		return nil, err
	}
	cfg.Database = database

	var cors Cors
	if err := env.Parse(&cors); err != nil {
		return nil, err
	}
	cfg.CORS = cors

	return &cfg, nil
}
