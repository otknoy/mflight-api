package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

// AppConfig has the configuration for entire application
type AppConfig struct {
	Port    int
	MfLight MfLightConfig
}

// Load loads the configuration
func Load() (AppConfig, error) {
	port, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		return AppConfig{}, errors.New("invalid port")
	}

	sc := AppConfig{
		port,
		loadMfLightConfig(),
	}

	return sc, nil
}

// MfLightConfig has the configuration to connect MfLight
type MfLightConfig struct {
	URL      string
	MobileID string
	CacheTTL time.Duration
}

func loadMfLightConfig() MfLightConfig {
	return MfLightConfig{
		URL:      os.Getenv("APP_MFLIGHT_URL"),
		MobileID: os.Getenv("APP_MFLIGHT_MOBILE_ID"),
		CacheTTL: parseDuration(os.Getenv("APP_MFLIGHT_CACHE_TTL")),
	}
}

func parseDuration(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}
