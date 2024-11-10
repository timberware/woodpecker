package config

import (
	"fmt"
	"os"
	"strconv"

	"woodpecker/internal/constants"
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

func LoadConfig() (*Config, error) {
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
		err = fmt.Errorf("missing required IP Service field in %s", constants.ConfigFilename)
		return nil, err
	}

	if config.PorkbunSecretKey == "" && config.NamecheapPassword == "" {
		err = fmt.Errorf("no domain variables specified in %s, exiting", constants.ConfigFilename)
		return nil, err
	}

	return config, nil
}
