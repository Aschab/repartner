package config

import (
	"encoding/json"
	"errors"
	"os"
)

// Config holds the application configuration.
type Config struct {
	PackSizes []int `json:"pack_sizes"`
}

// Load reads the configuration from the file specified by PACKS_CONFIG_PATH.
func Load() (*Config, error) {
	path := os.Getenv("PACKS_CONFIG_PATH")
	if path == "" {
		path = "configs/packs.json"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate checks that the configuration is valid.
func (c *Config) Validate() error {
	if len(c.PackSizes) == 0 {
		return errors.New("pack_sizes cannot be empty")
	}

	seen := make(map[int]bool)
	for _, size := range c.PackSizes {
		if size <= 0 {
			return errors.New("pack sizes must be positive integers")
		}
		if seen[size] {
			return errors.New("duplicate pack sizes are not allowed")
		}
		seen[size] = true
	}

	return nil
}
