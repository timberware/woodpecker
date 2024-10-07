package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	constants "woodpecker/config"
)

type Config struct {
	APIKey    string
	SecretKey string
	Domain    string
	Subdomain string
}

func LoadConfig(filename string) (*Config, error) {
	err := godotenv.Load(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file %s: %v", filename, err)
	}

	config := &Config{
		APIKey:    os.Getenv("PORKBUN_API_KEY"),
		SecretKey: os.Getenv("PORKBUN_SECRET_KEY"),
		Domain:    os.Getenv("DOMAIN"),
		Subdomain: os.Getenv("SUBDOMAIN"),
	}

	if config.APIKey == "" || config.SecretKey == "" || config.Domain == "" {
		return nil, fmt.Errorf("missing required fields in %s", filename)
	}

	return config, nil
}

func main() {
	config, err := LoadConfig(constants.Filename)
	if err != nil {
		fmt.Println("error loading env file:", err)
		os.Exit(-1)
	}

	fmt.Printf("config: %+v\n", config)
	fmt.Printf("config: %+v\n", constants.EditByNameTypeURL+config.Domain+"/A/"+config.Subdomain)
}
