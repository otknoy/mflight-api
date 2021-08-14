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

type MetricsGetter interface {
	GetMetrics(ctx context.Context) ([]Metrics, error)
}
