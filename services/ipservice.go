package services

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"woodpecker/config"
)

func GetPublicIP(config *config.Config) (string, error) {
	resp, err := http.Get(config.IPService)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve IP from %s: %v", config.IPService, err)
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response from %s: %v", config.IPService, err)
	}

	return strings.TrimSpace(string(ip)), nil
}
