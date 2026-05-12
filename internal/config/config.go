// Package config handles reading and writing persistent configuration.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// AppName is the application name used for config directory naming.
const AppName = "RADiscordPresence"

// OverrideDir, if non-empty, bypasses platform-specific directory resolution
// in Dir(). Set at startup when --config-dir is present in os.Args so the
// Windows service process (running as LocalSystem, with no user APPDATA) can
// locate the actual user's config file.
var OverrideDir string

// Config holds all persistent configuration for RADiscordPresence.
type Config struct {
	Username string `json:"username"`
	// APIKey is stored in the config file for now.
	// TODO: migrate to system keyring (Windows Credential Manager / macOS Keychain / libsecret).
	APIKey   string `json:"api_key"`
	Interval int    `json:"interval_seconds"`
}

// Default returns a Config with sensible defaults.
func Default() Config {
	return Config{Interval: 10}
}

// Dir returns the platform-appropriate config directory for this app.
// If OverrideDir is set it is returned directly; otherwise the platform
// default is derived from os.UserConfigDir().
func Dir() (string, error) {
	if OverrideDir != "" {
		return OverrideDir, nil
	}
	base, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine config directory: %w", err)
	}
	return filepath.Join(base, AppName), nil
}

// Load reads and returns the saved config. Returns Default() if no config file exists.
func Load() (Config, error) {
	dir, err := Dir()
	if err != nil {
		return Default(), err
	}
	data, err := os.ReadFile(filepath.Join(dir, "config.json"))
	if os.IsNotExist(err) {
		return Default(), nil
	}
	if err != nil {
		return Default(), fmt.Errorf("reading config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Default(), fmt.Errorf("parsing config: %w", err)
	}
	if cfg.Interval <= 0 {
		cfg.Interval = 10
	}
	return cfg, nil
}

// Save writes cfg to disk, creating the config directory if needed.
func Save(cfg Config) error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "config.json"), data, 0o600)
}
