package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-integrity-checker
type Config struct {
	ZebedeeRoot string `envconfig:"ZEBEDEE_ROOT"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		ZebedeeRoot: "content",
	}

	return cfg, envconfig.Process("", cfg)
}
