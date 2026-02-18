package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds application configuration
type Config struct {
	VaultDir         string   `json:"vault_dir"`
	Theme            string   `json:"theme"`
	RecentFiles      []string `json:"recent_files"`
	AutoSave         bool     `json:"auto_save"`
	AutoSaveInterval int      `json:"auto_save_interval"` // seconds
}

// DefaultConfig returns the default configuration
func DefaultConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting home directory: %w", err)
	}

	return &Config{
		VaultDir:         filepath.Join(homeDir, ".librenotes"),
		Theme:            "cyberpunk",
		RecentFiles:      []string{},
		AutoSave:         true,
		AutoSaveInterval: 30,
	}, nil
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "librenotes", "config.json"), nil
}

// ensureConfigDir ensures the config directory exists
func ensureConfigDir() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configDir := filepath.Join(homeDir, ".config", "librenotes")
	return os.MkdirAll(configDir, 0755)
}

// LoadConfig loads configuration from file or creates default
func LoadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return DefaultConfig()
	}

	// Try to load existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Config doesn't exist, create default
			cfg, err := DefaultConfig()
			if err != nil {
				return nil, err
			}
			// Save default config
			if err := cfg.Save(); err != nil {
				// Non-fatal: just log and continue
				fmt.Fprintf(os.Stderr, "Warning: could not save default config: %v\n", err)
			}
			return cfg, nil
		}
		return DefaultConfig()
	}

	// Parse existing config
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		// Invalid config, return default
		return DefaultConfig()
	}

	// Ensure vault directory is set
	if cfg.VaultDir == "" {
		homeDir, _ := os.UserHomeDir()
		cfg.VaultDir = filepath.Join(homeDir, ".librenotes")
	}

	return &cfg, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	if err := ensureConfigDir(); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
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

// GetRecentFiles returns the recent files list
func (c *Config) GetRecentFiles() []string {
	return c.RecentFiles
}

// SetRecentFiles sets the recent files list
func (c *Config) SetRecentFiles(files []string) {
	c.RecentFiles = files
}

// GetAutoSave returns whether autosave is enabled
func (c *Config) GetAutoSave() bool {
	return c.AutoSave
}

// GetAutoSaveInterval returns the autosave interval in seconds
func (c *Config) GetAutoSaveInterval() int {
	return c.AutoSaveInterval
}
