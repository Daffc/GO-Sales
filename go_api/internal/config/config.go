package config

import (
	"log"

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

type Server struct {
	Port               string `envconfig:"SERVER_PORT" default:"8080"`
	JwtSigningKey      []byte `envconfig:"JWT_SIGNING_KEY" required:"true"`
	JwtSessionDuration uint   `envconfig:"JWT_SESSION_DURATION" default:"24"`
}
type Config struct {
	Database Database
	Server   Server
}

func NewConfigParser() (*Config, error) {
	_ = godotenv.Load()
	cnf := Config{}
	err := envconfig.Process("", &cnf)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &cnf, nil
}
