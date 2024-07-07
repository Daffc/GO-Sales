package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Database struct {
	Name     string `envconfig:"DB_NAME" required:"true"`
	Host     string `envconfig:"DB_HOST" required:"true"`
	Port     string `envconfig:"DB_PORT" required:"true"`
	User     string `envconfig:"DB_USERNAME" required:"true"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
}

type Config struct {
	Database   Database
	ServerPort string `envconfig:"SERVER_PORT" default:"8080"`
}

func NewConfigParser() (*Config, error) {
	_ = godotenv.Load()
	cnf := Config{}
	err := envconfig.Process("", &cnf)
	return &cnf, err
}
