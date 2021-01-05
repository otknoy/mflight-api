package handler_test

import (
	"bytes"
	"errors"
	"mflight-exporter/domain"
	"mflight-exporter/handler"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubSensor struct{}

func (s *stubSensor) GetMetrics() (domain.Metrics, error) {
	return domain.Metrics{
		Temperature: domain.Temperature(20.0),
		Humidity:    domain.Humidity(50.0),
		Illuminance: domain.Illuminance(400),
	}, nil
}

type stubErrorSensor struct{}

func (s *stubErrorSensor) GetMetrics() (domain.Metrics, error) {
	return domain.Metrics{}, errors.New("failed to get metrics")
}

func TestServeHTTP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/getSensorMetrics", &bytes.Buffer{})
	got := httptest.NewRecorder()

	want := `{"temperature": 20.0, "humidity": 50.0, "illuminance": 400}`

	h := handler.NewSensorMetricsHandler(&stubSensor{})

	h.ServeHTTP(got, req)

	if v := got.Code; v != http.StatusOK {
		t.Errorf("http status: 200, but %v\n", v)
	}
	if v := got.Body.String(); v != want {
		t.Errorf("invalid response json:\nwant=%v\n got=%v\n", want, v)
	}
}

func TestServeHTTP_sensor_error(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/getSensorMetrics", &bytes.Buffer{})
	got := httptest.NewRecorder()

	h := handler.NewSensorMetricsHandler(&stubErrorSensor{})

	h.ServeHTTP(got, req)

	if v := got.Code; v != http.StatusInternalServerError {
		t.Errorf("http status: 500, but %v\n", v)
	}
	if v := got.Body.String(); v != "{}" {
		t.Errorf("empty response:\nwant={}\n got=%v\n", v)
	}
}
