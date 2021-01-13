package application

import (
	"context"
	"mflight-api/domain"
)

// MetricsCollector is interface to collect metrics
type MetricsCollector interface {
	CollectMetrics(ctx context.Context) (domain.Metrics, error)
}

// NewMetricsCollector creates a new MetricsCollector Based on domain.Sensor
func NewMetricsCollector(s domain.Sensor) MetricsCollector {
	return &metricsCollector{s}
}

type metricsCollector struct {
	sensor domain.Sensor
}

// CollectMetrics returns collected metrics
func (c *metricsCollector) CollectMetrics(ctx context.Context) (domain.Metrics, error) {
	return c.sensor.GetMetrics(ctx)
}
