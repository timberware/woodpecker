package services

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetPublicIP() (string, error) {
	services := []string{
		"https://api.ipify.org",
		"https://ifconfig.me",
		"https://ipinfo.io/ip",
	}

	for _, service := range services {
		resp, err := http.Get(service)
		if err != nil {
			fmt.Printf("failed to retrieve IP from %s: %v\n", service, err)
			continue
		}
		defer resp.Body.Close()

		ip, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("failed to read response from %s: %v\n", service, err)
			continue
		}

		return strings.TrimSpace(string(ip)), nil
	}

	return "", fmt.Errorf("cannot obtain IP address: all services failed or no network connection")
}
