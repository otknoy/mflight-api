package handler

import (
	"encoding/json"
	"mflight-api/domain"
	"net/http"

	"go.uber.org/zap"
)

// SensorMetricsHandler is struct to get sensor metrics
type SensorMetricsHandler struct {
	metricsGetter domain.MetricsGetter
}

type response []metrics

type metrics struct {
	Unixtime    int64   `json:"unixtime"`
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
	Illuminance int16   `json:"illuminance"`
}

// NewSensorMetricsHandler creates a new SensorMetricsHandler based on domain.MetricsGetter
func NewSensorMetricsHandler(g domain.MetricsGetter) *SensorMetricsHandler {
	return &SensorMetricsHandler{g}
}

// ServeHTTP implements http.Handler
func (h *SensorMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l, err := h.metricsGetter.GetMetrics(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := convert(l)

	bytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	successResponse(w, bytes)
}

func convert(l domain.MetricsList) response {
	res := make([]metrics, len(l))
	for i, m := range l {
		res[i] = metrics{
			Unixtime:    m.Time.Unix(),
			Temperature: float32(m.Temperature),
			Humidity:    float32(m.Humidity),
			Illuminance: int16(m.Illuminance),
		}
	}
	return res
}

func successResponse(w http.ResponseWriter, bytes []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(bytes)
	if err != nil {
		zap.L().Error("Write failed", zap.Error(err))
	}
}
