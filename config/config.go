package config

import (
	"errors"
	"os"
	"strconv"
)

// AppConfig has the configuration for entire application
type AppConfig struct {
	Port    int
	MfLight MfLightConfig
}

// MfLightConfig has the configuration to connect MfLight
type MfLightConfig struct {
	URL      string
	MobileID string
}

// Load loads the configuration
func Load() (AppConfig, error) {
	port, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		return AppConfig{}, errors.New("invalid port")
	}

	sc := AppConfig{
		port,
		MfLightConfig{
			os.Getenv("APP_MFLIGHT_URL"),
			os.Getenv("APP_MFLIGHT_MOBILE_ID"),
		},
	}

	return sc, nil
}
