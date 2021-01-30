package mflight

import (
	"context"
	"fmt"
	"mflight-api/domain"
)

type mfLightSensor struct {
	client Client
}

// NewMfLightSensor creates a new MfLight based on mflight.Client
func NewMfLightSensor(c Client) domain.Sensor {
	return &mfLightSensor{c}
}

// GetMetrics returns current Metrics
func (l *mfLightSensor) GetMetrics(ctx context.Context) (domain.TimeSeriesMetrics, error) {
	res, err := l.client.GetSensorMonitor(ctx)
	if err != nil {
		return domain.TimeSeriesMetrics{}, err
	}

	last := len(res.Tables) - 1
	if last < 0 {
		return domain.TimeSeriesMetrics{}, fmt.Errorf("invalid api response: %v", res)
	}

	table := res.Tables[len(res.Tables)-1]

	m := domain.Metrics{
		Temperature: domain.Temperature(table.Temperature),
		Humidity:    domain.Humidity(table.Humidity),
		Illuminance: domain.Illuminance(table.Illuminance),
	}

	return domain.TimeSeriesMetrics([]domain.Metrics{m}), nil
}
