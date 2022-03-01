package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// AppConfig has the configuration for entire application
type AppConfig struct {
	Port    int
	MfLight MfLightConfig
}

// MfLightConfig has the configuration to connect MfLight
type MfLightConfig struct {
	URL      string
	MobileID string        `split_words:"true"`
	CacheTTL time.Duration `split_words:"true"`
}

// Load loads the configuration
func Load() (AppConfig, error) {
	var c AppConfig
	err := envconfig.Process("app", &c)
	if err != nil {
		return AppConfig{}, fmt.Errorf("invalid config: %w", err)
	}

	return c, nil
}
