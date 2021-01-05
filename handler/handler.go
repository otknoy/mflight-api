package handler

import (
	"fmt"
	"mflight-exporter/domain"
	"net/http"
)

// SensorMetricsHandler is struct to get sensor metrics
type SensorMetricsHandler struct {
	sensor domain.Sensor
}

// NewSensorMetricsHandler creates a new SensorMetricsHandler based on domain.Sensor
func NewSensorMetricsHandler(sensor domain.Sensor) *SensorMetricsHandler {
	return &SensorMetricsHandler{sensor}
}

// ServeHTTP implements http.Handler
func (h *SensorMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m, err := h.sensor.GetMetrics()

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "{}")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(
		w,
		"{\"temperature\": %0.1f, \"humidity\": %0.1f, \"illuminance\": %d}",
		m.Temperature, m.Humidity, m.Illuminance,
	)
}
