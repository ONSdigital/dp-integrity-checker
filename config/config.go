package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-integrity-checker
type Config struct {
	ZebedeeRoot                string `envconfig:"ZEBEDEE_ROOT"`
	CheckPublishedPreviousDays int    `envconfig:"CHECK_PUBLISHED_PREVIOUS_DAYS"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		ZebedeeRoot:                "content",
		CheckPublishedPreviousDays: 1,
	}

	return cfg, envconfig.Process("", cfg)
}
