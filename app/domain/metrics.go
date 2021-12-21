package domain

import (
	"context"
	"errors"
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

type MetricsList []Metrics

func (l MetricsList) Last() (Metrics, error) {
	if len(l) == 0 {
		return Metrics{}, errors.New("empty")
	}

	last := len(l) - 1

	return l[last], nil
}

type MetricsGetter interface {
	GetMetrics(ctx context.Context) (MetricsList, error)
}
