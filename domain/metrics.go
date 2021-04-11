package domain

import (
	"context"
	"time"
)

// Temperature is value object
type Temperature float32

// Humidity is value object
type Humidity float32

// Illuminance is value object
type Illuminance int16

// Metrics has multiple sensor values
type Metrics struct {
	Time        time.Time
	Temperature Temperature
	Humidity    Humidity
	Illuminance Illuminance
}

// TimeSeriesMetrics is metrics list in time series order.
type TimeSeriesMetrics []Metrics

// MetricsRepository is interface to get metrics
type MetricsRepository interface {
	GetMetrics(ctx context.Context) (TimeSeriesMetrics, error)
}
