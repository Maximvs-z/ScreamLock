// Package config handles loading and saving ScreamLock configuration.
// Config is stored in JSON for simplicity and human editability.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds application settings. ThresholdDB is in decibels (e.g. -50 = loud).
// DeviceID should match one of the IDs from -list-devices output.
type Config struct {
	// DeviceID is the Windows audio device ID (from enumeration). Empty = use default capture device.
	DeviceID string `json:"device_id"`
	// ThresholdDB is the level in dB above which the workstation is locked. Typical: -50 to -30.
	ThresholdDB float64 `json:"threshold_db"`
	// CheckIntervalSeconds is how often to sample the microphone (default 1).
	CheckIntervalSeconds int `json:"check_interval_seconds"`
}

// DefaultConfig returns sensible defaults matching the original script (-50 dB, 1 second).
func DefaultConfig() Config {
	return Config{
		DeviceID:             "",
		ThresholdDB:          -50,
		CheckIntervalSeconds: 1,
	}
}

// ConfigDir returns the application data directory for ScreamLock (e.g. %APPDATA%\ScreamLock).
func ConfigDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "ScreamLock"), nil
}

// ConfigPath returns the full path to config.json.
func ConfigPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// Load reads config from the standard path. Creates default config file if missing.
func Load() (Config, string, error) {
	path, err := ConfigPath()
	if err != nil {
		return Config{}, "", err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := DefaultConfig()
			if err := Save(cfg); err != nil {
				return cfg, path, err
			}
			return cfg, path, nil
		}
		return Config{}, path, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, path, err
	}
	if cfg.CheckIntervalSeconds <= 0 {
		cfg.CheckIntervalSeconds = 1
	}
	return cfg, path, nil
}

// Save writes config to the standard path, creating the directory if needed.
func Save(cfg Config) error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := filepath.Join(dir, "config.json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
