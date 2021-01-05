package handler

import (
	"fmt"
	"mflight-exporter/application"
	"net/http"
)

// SensorMetricsHandler is struct to get sensor metrics
type SensorMetricsHandler struct {
	metricsCollector application.MetricsCollector
}

// NewSensorMetricsHandler creates a new SensorMetricsHandler based on domain.Sensor
func NewSensorMetricsHandler(c application.MetricsCollector) *SensorMetricsHandler {
	return &SensorMetricsHandler{c}
}

// ServeHTTP implements http.Handler
func (h *SensorMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m, err := h.metricsCollector.CollectMetrics()

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
