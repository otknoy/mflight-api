package mflight_test

import (
	"mflight-exporter/infrastructure/mflight"
	"testing"
)

func TestGetMetrics(t *testing.T) {
	s := NewStubServer(t)
	defer s.Close()

	sensor := mflight.NewMfLightSensor(s.URL, "test-mobile-id")

	m, err := sensor.GetMetrics()
	if err != nil {
		t.Fatal(err)
	}

	if v := m.Temperature; v != 21.9 {
		t.Errorf("invalid temperature: %v\n", v)
	}
	if v := m.Humidity; v != 43.0 {
		t.Errorf("invalid humidity: %v\n", v)
	}
	if v := m.Illuminance; v != 406 {
		t.Errorf("invalid illuminance: %v\n", v)
	}
}
