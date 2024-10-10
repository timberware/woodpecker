package config

import (
	"fmt"
	"os"
	"woodpecker/constants"

	"github.com/joho/godotenv"
)

type Config struct {
	APIKey    string
	SecretKey string
	Domain    string
	Subdomain string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(constants.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file %s: %v", constants.Filename, err)
	}

	config := &Config{
		APIKey:    os.Getenv("PORKBUN_API_KEY"),
		SecretKey: os.Getenv("PORKBUN_SECRET_KEY"),
		Domain:    os.Getenv("DOMAIN"),
		Subdomain: os.Getenv("SUBDOMAIN"),
	}

	if config.APIKey == "" || config.SecretKey == "" || config.Domain == "" {
		return nil, fmt.Errorf("missing required fields in %s", constants.Filename)
	}

	return config, nil
}
