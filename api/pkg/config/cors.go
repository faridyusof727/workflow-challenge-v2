package config

type Cors struct {
	AllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS"`
	AllowedMethods []string `env:"CORS_ALLOWED_METHODS"`
	AllowedHeaders []string `env:"CORS_ALLOWED_HEADERS"`
}
