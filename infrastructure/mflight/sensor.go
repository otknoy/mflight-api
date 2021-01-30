package mflight

import (
	"context"
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

	return convert(res.Tables), nil
}

func convert(tables []Table) domain.TimeSeriesMetrics {
	ts := make([]domain.Metrics, len(tables))
	for i, t := range tables {
		ts[i] = domain.Metrics{
			Temperature: domain.Temperature(t.Temperature),
			Humidity:    domain.Humidity(t.Humidity),
			Illuminance: domain.Illuminance(t.Illuminance),
		}
	}
	return domain.TimeSeriesMetrics(ts)
}
