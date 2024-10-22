package config

import (
	"fmt"
	"os"
	"woodpecker/constants"

	"github.com/joho/godotenv"
)

type Config struct {
	APIKey                string
	SecretKey             string
	Domain                string
	Subdomain             string
	IPService             string
	EditByNameTypeURL     string
	RetrieveByNameTypeURL string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(constants.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file %s: %v", constants.Filename, err)
	}

	config := &Config{
		APIKey:                os.Getenv("PORKBUN_API_KEY"),
		SecretKey:             os.Getenv("PORKBUN_SECRET_KEY"),
		Domain:                os.Getenv("DOMAIN"),
		Subdomain:             os.Getenv("SUBDOMAIN"),
		IPService:             os.Getenv("IP_SERVICE"),
		EditByNameTypeURL:     os.Getenv("PORKBUN_API_EDIT_URL"),
		RetrieveByNameTypeURL: os.Getenv("PORKBUN_API_RETRIEVE_URL"),
	}

	if config.APIKey == "" || config.SecretKey == "" {
		return nil, fmt.Errorf("missing required key fields in %s", constants.Filename)
	}

	if config.Domain == "" {
		return nil, fmt.Errorf("missing required domain field in %s", constants.Filename)
	}

	if config.EditByNameTypeURL == "" || config.RetrieveByNameTypeURL == "" {
		return nil, fmt.Errorf("missing required Porkbun URL fields in %s", constants.Filename)
	}

	if config.IPService == "" {
		return nil, fmt.Errorf("missing required IP Service field in %s", constants.Filename)
	}

	return config, nil
}
