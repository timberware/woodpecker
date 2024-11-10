package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"woodpecker/internal/constants"
)

func GetAppPath() (string, error) {
	dir := constants.CurrentIPFilenameDir

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	return dir, nil
}

func ReadIPFromFile(configPath string) (string, error) {
	ipFile := filepath.Join(configPath, constants.CurrentIpFilename)
	data, err := os.ReadFile(ipFile)
	if os.IsNotExist(err) {
		Log.Info().Msg("IP file does not exist yet, assuming first run")
		return "", nil
	}

	if err != nil {
		return "", fmt.Errorf("failed to read IP from file: %w", err)
	}

	return string(data), nil
}

func WriteIPToFile(ip, configPath string) error {
	ipFile := filepath.Join(configPath, constants.CurrentIpFilename)
	err := os.WriteFile(ipFile, []byte(ip), 0644)
	if err != nil {
		return fmt.Errorf("failed to write IP to file: %w", err)
	}

	Log.Info().Str("level", "update").Msgf("stored IP address (%s) locally", ip)
	return nil
}
