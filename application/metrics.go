package application

import (
	"context"
	"errors"
	"mflight-api/domain"
)

// MetricsCollector is interface to collect metrics
type MetricsCollector interface {
	CollectLatestMetrics(ctx context.Context) (domain.Metrics, error)
	CollectMetricsList(ctx context.Context) ([]domain.Metrics, error)
}

// NewMetricsCollector creates a new MetricsCollector based on domain.MetricsRepository
func NewMetricsCollector(g domain.MetricsGetter) MetricsCollector {
	return &metricsCollector{g}
}

type metricsCollector struct {
	g domain.MetricsGetter
}

// CollectLatestMetrics returns collected metrics
func (c *metricsCollector) CollectLatestMetrics(ctx context.Context) (domain.Metrics, error) {
	ts, err := c.g.GetMetrics(ctx)
	if err != nil {
		return domain.Metrics{}, err
	}

	last := len(ts) - 1
	if last < 0 {
		return domain.Metrics{}, errors.New("no metrics")
	}

	return ts[last], nil
}

// CollectMetricsList returns collected metrics list
func (c *metricsCollector) CollectMetricsList(ctx context.Context) ([]domain.Metrics, error) {
	return c.g.GetMetrics(ctx)
}
