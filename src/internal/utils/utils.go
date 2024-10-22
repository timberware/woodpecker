package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"woodpecker/src/constants"
)

func GetAppPath() (string, error) {
	var dir string

	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		dir = filepath.Join(appData, "woodpecker")
	case "linux":
		home := os.Getenv("HOME")
		dir = filepath.Join(home, ".local", "share", "woodpecker")
	}

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	return dir, nil
}

func ReadIPFromFile(configPath string) (string, error) {
	ipFile := filepath.Join(configPath, constants.CurrentIpFilename)
	data, err := os.ReadFile(ipFile)
	if os.IsNotExist(err) {
		fmt.Println("IP file does not exist yet, assuming first run")
		return "", nil
	}

	if err != nil {
		return "", fmt.Errorf("failed to read IP from file: %v", err)
	}

	return string(data), nil
}

func WriteIPToFile(ip, configPath string) error {
	ipFile := filepath.Join(configPath, constants.CurrentIpFilename)
	err := os.WriteFile(ipFile, []byte(ip), 0644)
	if err != nil {
		return fmt.Errorf("failed to write IP to file: %v", err)
	}

	return nil
}
