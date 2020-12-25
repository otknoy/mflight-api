package mflight

import (
	"fmt"
)

// Temperature is value object
type Temperature float32

// Humidity is value object
type Humidity float32

// Illuminance is value object
type Illuminance int16

// Metrics has multiple sensor values
type Metrics struct {
	Temperature Temperature `xml:"temp"`
	Humidity    Humidity    `xml:"humi"`
	Illuminance Illuminance `xml:"illu"`
}

// Sensor is interface to get metrics
type Sensor interface {
	GetMetrics() (Metrics, error)
}

type mfLightSensor struct {
	serverURL string
	mobileID  string
}

// NewMfLight creates a new MfLight based on mflight server configuration
func NewMfLight(serverURL, mobileID string) Sensor {
	return &mfLightSensor{serverURL, mobileID}
}

// GetMetrics returns current Metrics
func (l *mfLightSensor) GetMetrics() (Metrics, error) {
	res, err := getSensorMonitor(l.serverURL, l.mobileID)
	if err != nil {
		return Metrics{}, err
	}

	tables := res.Tables
	last := len(tables) - 1
	if last < 0 {
		return Metrics{}, fmt.Errorf("invalid api response: %v", res)
	}

	table := tables[last]

	m := Metrics{
		Temperature(table.Temperature),
		Humidity(table.Humidity),
		Illuminance(table.Illuminance),
	}

	return m, nil
}
