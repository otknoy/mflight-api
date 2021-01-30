package domain

import "context"

// Temperature is value object
type Temperature float32

// Humidity is value object
type Humidity float32

// Illuminance is value object
type Illuminance int16

// Metrics has multiple sensor values
type Metrics struct {
	Temperature Temperature
	Humidity    Humidity
	Illuminance Illuminance
}

// TimeSeriesMetrics is metrics list in time series order.
type TimeSeriesMetrics []Metrics

// Sensor is interface to get metrics
type Sensor interface {
	GetMetrics(ctx context.Context) (TimeSeriesMetrics, error)
}
