package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-integrity-checker
type Config struct {
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{}

	return cfg, envconfig.Process("", cfg)
}
