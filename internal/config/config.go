package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds application configuration
type Config struct {
	VaultDir string
	Theme    string
}

// DefaultConfig returns the default configuration
func DefaultConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting home directory: %w", err)
	}

	return &Config{
		VaultDir: filepath.Join(homeDir, ".librenotes"),
		Theme:    "cyberpunk",
	}, nil
}

// LoadConfig loads configuration from file or creates default
func LoadConfig() (*Config, error) {
	// For now, just return default config
	// Later we can add YAML/TOML file loading
	return DefaultConfig()
}

// GetVaultDir returns the vault directory
func (c *Config) GetVaultDir() string {
	return c.VaultDir
}

// GetTheme returns the current theme name
func (c *Config) GetTheme() string {
	return c.Theme
}

// SetTheme sets the theme
func (c *Config) SetTheme(theme string) {
	c.Theme = theme
}
