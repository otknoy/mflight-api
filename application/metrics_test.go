package application_test

import (
	"mflight-exporter/application"
	"mflight-exporter/domain"
	"testing"
)

type stubSensor struct{}

func (s *stubSensor) GetMetrics() (domain.Metrics, error) {
	return domain.Metrics{
		Temperature: domain.Temperature(18.0),
		Humidity:    domain.Humidity(45.0),
		Illuminance: domain.Illuminance(300),
	}, nil
}

func TestCollectMetrics(t *testing.T) {
	c := application.NewMetricsCollector(&stubSensor{})

	m, _ := c.CollectMetrics()

	if m.Temperature != 18.0 ||
		m.Humidity != 45.0 ||
		m.Illuminance != 300 {

		t.Errorf("returned metrics is invalid: %v\n", m)
	}
}
