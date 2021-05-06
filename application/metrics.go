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

// NewMetricsCollector creates a new MetricsCollector based on domain.MetricsRepository
func NewMetricsCollector(s domain.MetricsRepository) MetricsCollector {
	return &metricsCollector{s}
}

type metricsCollector struct {
	sensor domain.MetricsRepository
}

// CollectLatestMetrics returns collected metrics
func (c *metricsCollector) CollectLatestMetrics(ctx context.Context) (domain.Metrics, error) {
	ts, err := c.sensor.GetMetrics(ctx)
	if err != nil {
		return domain.Metrics{}, err
	}

	last := len(ts) - 1
	if last < 0 {
		return domain.Metrics{}, errors.New("no metrics")
	}

	return ts[last], nil
}

// CollectTimeSeriesMetrics returns collected metrics list
func (c *metricsCollector) CollectTimeSeriesMetrics(ctx context.Context) (domain.TimeSeriesMetrics, error) {
	return c.sensor.GetMetrics(ctx)
}
