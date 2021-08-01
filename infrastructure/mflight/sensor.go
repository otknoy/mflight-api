package mflight

import (
	"context"
	"mflight-api/domain"
	"mflight-api/infrastructure/mflight/httpclient"
	"time"
)

type mfLightSensor struct {
	client httpclient.Client
}

// NewMfLightSensor creates a new domain.MetricsRepository based on mflight.Client
func NewMfLightSensor(c httpclient.Client) domain.MetricsRepository {
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

func convert(tables []httpclient.Table) domain.TimeSeriesMetrics {
	ts := make([]domain.Metrics, len(tables))
	for i, t := range tables {
		ts[i] = domain.Metrics{
			Time:        time.Unix(t.Unixtime, 0),
			Temperature: domain.Temperature(t.Temperature),
			Humidity:    domain.Humidity(t.Humidity),
			Illuminance: domain.Illuminance(t.Illuminance),
		}
	}
	return ts
}
