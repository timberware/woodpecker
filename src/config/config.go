package config

import (
	"fmt"
	"os"
	"strconv"
	"woodpecker/src/constants"

	"github.com/joho/godotenv"
)

type Config struct {
	IPService                    string
	CheckInterval                int
	PorkbunAPIKey                string
	PorkbunSecretKey             string
	PorkbunDomain                string
	PorkbunSubdomain             string
	PorkbunEditByNameTypeURL     string
	PorkbunRetrieveByNameTypeURL string
	NamecheapEditURL             string
	NamecheapPassword            string
	NamecheapDomain              string
	NamecheapSubdomain           string
}

func LoadConfig(configPath string) (*Config, error) {
	err := godotenv.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file %s: %v", constants.ConfigFilename, err)
	}

	intervalStr := os.Getenv("CHECK_INTERVAL")
	interval, err := strconv.Atoi(intervalStr)
	if err != nil || interval <= 0 {
		interval = 3
	}

	config := &Config{
		IPService:                    os.Getenv("IP_SERVICE"),
		CheckInterval:                interval,
		PorkbunAPIKey:                os.Getenv("PORKBUN_API_KEY"),
		PorkbunSecretKey:             os.Getenv("PORKBUN_SECRET_KEY"),
		PorkbunDomain:                os.Getenv("PORKBUN_DOMAIN"),
		PorkbunSubdomain:             os.Getenv("PORKBUN_SUBDOMAIN"),
		PorkbunEditByNameTypeURL:     os.Getenv("PORKBUN_API_EDIT_URL"),
		PorkbunRetrieveByNameTypeURL: os.Getenv("PORKBUN_API_RETRIEVE_URL"),
		NamecheapEditURL:             os.Getenv("NAMECHEAP_EDIT_URL"),
		NamecheapPassword:            os.Getenv("NAMECHEAP_PASSWORD"),
		NamecheapDomain:              os.Getenv("NAMECHEAP_DOMAIN"),
		NamecheapSubdomain:           os.Getenv("NAMECHEAP_SUBDOMAIN"),
	}

	if config.IPService == "" {
		return nil, fmt.Errorf("missing required IP Service field in %s", constants.ConfigFilename)
	}

	if config.PorkbunAPIKey == "" || config.PorkbunSecretKey == "" || config.PorkbunDomain == "" || config.PorkbunEditByNameTypeURL == "" || config.PorkbunRetrieveByNameTypeURL == "" {
		return nil, fmt.Errorf("missing required Porkbun key fields in %s", constants.ConfigFilename)
	}

	if config.NamecheapPassword == "" || config.NamecheapDomain == "" || config.NamecheapEditURL == "" {
		return nil, fmt.Errorf("missing required Namecheap key fields in %s", constants.ConfigFilename)
	}

	return config, nil
}
