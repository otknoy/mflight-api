package main

import (
	"fmt"
	"log"
	"mflight-api/application"
	"mflight-api/config"
	"mflight-api/handler"
	"mflight-api/infrastructure/mflight"
	"mflight-api/infrastructure/prometheus/collector"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	c, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	sensor := mflight.NewMfLightSensor(c.MfLight.URL, c.MfLight.MobileID)
	metricsCollector := application.NewMetricsCollector(sensor)

	h := handler.NewSensorMetricsHandler(metricsCollector)
	http.Handle("/getSensorMetrics", h)

	col := collector.NewMfLightCollector(metricsCollector)
	prometheus.MustRegister(col)
	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", c.Port), nil))
}
