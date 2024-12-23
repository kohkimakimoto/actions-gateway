package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"strings"
)

// Config is the configuration for the server
type Config struct {
	// Addr is the address to listen on
	Addr string `toml:"addr"`
	// URL is the base URL of the server
	URL string `toml:"url"`
	// Secret is a HS256 secret key. It must be longer than 256 bits.
	Secret string `toml:"secret"`
	// ExposeNewToken enables the /new-token and /api/new-token endpoints
	ExposeNewToken bool `toml:"expose_new_token"`
	// Debug enables debug logging
	Debug bool `toml:"debug"`
}

func New() *Config {
	return &Config{
		Addr:   ":18800",
		URL:    "http://localhost:18800",
		Secret: "",
	}
}

// WebSocketURL returns the WebSocket URL based on the server URL configuration
func (c *Config) WebSocketURL() string {
	var wsPath string
	if strings.HasPrefix(c.URL, "https://") {
		wsPath = "wss://" + c.URL[8:]
	} else if strings.HasPrefix(c.URL, "http://") {
		wsPath = "ws://" + c.URL[7:]
	} else {
		wsPath = "ws://" + c.URL
	}
	return strings.TrimRight(wsPath, "/")
}

// UpdateByFile updates the configuration from a file
func UpdateByFile(c *Config, path string) error {
	if _, err := toml.DecodeFile(path, c); err != nil {
		return fmt.Errorf("failed to load config from file: %w", err)
	}
	return nil
}

// UpdateByEnvironments updates the configuration from environment variables
func UpdateByEnvironments(c *Config) {
	if v := os.Getenv("ACTIONS_GATEWAY_ADDR"); v != "" {
		c.Addr = v
	}
	if v := os.Getenv("ACTIONS_GATEWAY_URL"); v != "" {
		c.URL = v
	}
	if v := os.Getenv("ACTIONS_GATEWAY_SECRET"); v != "" {
		c.Secret = v
	}
	if v := os.Getenv("ACTIONS_GATEWAY_EXPOSE_NEW_TOKEN"); v != "" {
		v = strings.ToLower(v)
		if v == "true" || v == "1" {
			c.ExposeNewToken = true
		} else {
			c.ExposeNewToken = false
		}
	}
	if v := os.Getenv("ACTIONS_GATEWAY_DEBUG"); v != "" {
		v = strings.ToLower(v)
		if v == "true" || v == "1" {
			c.Debug = true
		} else {
			c.Debug = false
		}
	}
}
