package config

import (
	"errors"
	"os"
	"strconv"
)

type ServerConfig struct {
	Port    int
	MfLight MfLightConfig
}

type MfLightConfig struct {
	URL      string
	MobileID string
}

func Load() (ServerConfig, error) {
	port, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		return ServerConfig{}, errors.New("invalid port")
	}

	sc := ServerConfig{
		port,
		MfLightConfig{
			os.Getenv("APP_MFLIGHT_URL"),
			os.Getenv("APP_MFLIGHT_MOBILE_ID"),
		},
	}

	return sc, nil
}
