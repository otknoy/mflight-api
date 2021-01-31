package application

import (
	"context"
	"errors"
	"mflight-api/domain"
)

// MetricsCollector is interface to collect metrics
type MetricsCollector interface {
	CollectLatestMetrics(ctx context.Context) (domain.Metrics, error)
	CollectTimeSeriesMetrics(ctx context.Context) (domain.TimeSeriesMetrics, error)
}

// NewMetricsCollector creates a new MetricsCollector Based on domain.Sensor
func NewMetricsCollector(s domain.Sensor) MetricsCollector {
	return &metricsCollector{s}
}

type metricsCollector struct {
	sensor domain.Sensor
}

// CollectLatestMetrics returns collected metrics
func (c *metricsCollector) CollectLatestMetrics(ctx context.Context) (domain.Metrics, error) {
	ts, _ := c.sensor.GetMetrics(ctx)

	last := len(ts) - 1
	if last < 0 {
		return domain.Metrics{}, errors.New("empty metrics")
	}

	return ts[last], nil
}

// CollectTimeSeriesMetrics returns collected metrics list
func (c *metricsCollector) CollectTimeSeriesMetrics(ctx context.Context) (domain.TimeSeriesMetrics, error) {
	return c.sensor.GetMetrics(ctx)
}
