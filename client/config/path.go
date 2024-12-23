package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// DefaultDir returns the default config directory path.
func DefaultDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get the user's home directory: %w", err)
	}

	// Check the XDG_CONFIG_HOME environment variable
	configHomeDir := os.Getenv("XDG_CONFIG_HOME")
	if configHomeDir == "" {
		// If the XDG_CONFIG_HOME environment variable is not set, use the $HOME/.config directory.
		configHomeDir = filepath.Join(homeDir, ".config")
	}

	return filepath.Join(configHomeDir, "actions-gateway"), nil
}

// DefaultFilepath returns the default config file path for the actions-gateway client.
func DefaultFilepath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Check the XDG_CONFIG_HOME environment variable
	configHomeDir := os.Getenv("XDG_CONFIG_HOME")
	if configHomeDir == "" {
		// If the XDG_CONFIG_HOME environment variable is not set, use the $HOME/.config directory.
		configHomeDir = filepath.Join(homeDir, ".config")
	}

	configDir := filepath.Join(configHomeDir, "actions-gateway")
	configFile := filepath.Join(configDir, "config.toml")

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// If the config file does not exist, try to find the config file in the $HOME/.actions-gateway directory.
		configDir = filepath.Join(homeDir, ".actions-gateway")
		configFile = filepath.Join(configDir, "config.toml")
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			return "", err
		}
	}

	return configFile, nil
}
