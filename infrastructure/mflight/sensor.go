package mflight

import (
	"fmt"
	"mflight-api/domain"
)

type mfLightSensor struct {
	client Client
}

// NewMfLightSensor creates a new MfLight based on mflight server configuration
func NewMfLightSensor(serverURL, mobileID string) domain.Sensor {
	return &mfLightSensor{
		NewClient(serverURL, mobileID),
	}
}

// GetMetrics returns current Metrics
func (l *mfLightSensor) GetMetrics() (domain.Metrics, error) {
	res, err := l.client.GetSensorMonitor()
	if err != nil {
		return domain.Metrics{}, err
	}

	tables := res.Tables
	last := len(tables) - 1
	if last < 0 {
		return domain.Metrics{}, fmt.Errorf("invalid api response: %v", res)
	}

	table := tables[last]

	m := domain.Metrics{
		Temperature: domain.Temperature(table.Temperature),
		Humidity:    domain.Humidity(table.Humidity),
		Illuminance: domain.Illuminance(table.Illuminance),
	}

	return m, nil
}
