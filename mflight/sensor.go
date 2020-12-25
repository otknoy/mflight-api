package mflight

import (
	"errors"
	"fmt"
)

type Temperature float32
type Humidity float32
type Illuminance int16

type Metrics struct {
	Temperature Temperature `xml:"temp"`
	Humidity    Humidity    `xml:"humi"`
	Illuminance Illuminance `xml:"illu"`
}

type Sensor interface {
	GetMetrics() (Metrics, error)
}

type mfLightSensor struct {
	serverURL string
	mobileID  string
}

func NewMfLight(serverURL, mobileID string) Sensor {
	return &mfLightSensor{serverURL, mobileID}
}

func (l *mfLightSensor) GetMetrics() (Metrics, error) {
	res, err := getSensorMonitor(l.serverURL, l.mobileID)
	if err != nil {
		return Metrics{}, err
	}

	tables := res.Tables
	last := len(tables) - 1
	if last < 0 {
		return Metrics{}, errors.New(fmt.Sprintf("invalid api response: %v", res))
	}

	table := tables[last]

	m := Metrics{
		Temperature(table.Temperature),
		Humidity(table.Humidity),
		Illuminance(table.Illuminance),
	}

	return m, nil
}
