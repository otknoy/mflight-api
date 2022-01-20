package server

import (
	"fmt"
	"mflight-api/app/domain"
	"mflight-api/app/handler"
	"mflight-api/app/infrastructure/prometheus/collector"
	"mflight-api/app/infrastructure/prometheus/middleware"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(metricsGetter domain.MetricsGetter, port int) GracefulShutdownServer {
	h := handler.NewSensorMetricsHandler(metricsGetter)

	prometheus.MustRegister(
		collector.NewMfLightCollector(metricsGetter),
	)

	mux := http.NewServeMux()
	mux.Handle("/getSensorMetrics", middleware.InstrumentHandlerMetrics(h))
	mux.Handle("/metrics", middleware.InstrumentHandlerMetrics(promhttp.Handler()))

	return GracefulShutdownServer{
		Server: http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}

}
