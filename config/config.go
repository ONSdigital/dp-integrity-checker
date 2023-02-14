package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-integrity-checker
type Config struct {
	ZebedeeRoot                string `envconfig:"ZEBEDEE_ROOT"`
	CheckPublishedPreviousDays int    `envconfig:"CHECK_PUBLISHED_PREVIOUS_DAYS"`
	SlackEnabled               bool   `envconfig:"SLACK_ENABLED"`
	SlackConfig                Slack
}

type Slack struct {
	ApiToken     string `envconfig:"SLACK_API_TOKEN"  json:"-"`
	UserName     string `envconfig:"SLACK_USER_NAME"`
	AlarmChannel string `envconfig:"SLACK_ALARM_CHANNEL"`
	AlarmEmoji   string `envconfig:"SLACK_ALARM_EMOJI"`
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
		SlackEnabled:               false,
		SlackConfig: Slack{
			UserName:     "Integrity Checker",
			AlarmChannel: "#sandbox-alarm",
			AlarmEmoji:   ":rotating_light:",
		},
	}

	return cfg, envconfig.Process("", cfg)
}
