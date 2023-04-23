// Package config handles environment variables.
package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

// Config contains environment variables.
type Config struct {
	Port                   string `envconfig:"PORT" default:"8080"`
	GoogleClientID         string `envconfig:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret     string `envconfig:"GOOGLE_CLIENT_SECRET"`
	GoogleOauthCallbackURL string `envconfig:"GOOGLE_OAUTH_CALLBACK_URL"`
	JWTSecret              string `envconfig:"JWT_SECRET"`
}

// LoadConfig reads environment variables and populates Config.
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Info("No .env file found")
	}

	var c Config

	err := envconfig.Process("", &c)

	return &c, err
}
