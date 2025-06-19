package config

type Database struct {
	URL string `env:"DATABASE_URL,required"`
}
