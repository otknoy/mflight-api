package domain

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

// Sensor is interface to get metrics
type Sensor interface {
	GetMetrics() (Metrics, error)
}
