package main

import (
	"fmt"
	"log"
	"mflight-exporter/collector"
	"mflight-exporter/config"
	"mflight-exporter/mflight"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	c, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	sensor := mflight.NewMfLight(c.MfLight.URL, c.MfLight.MobileID)
	col := collector.NewMfLightCollector(sensor)

	prometheus.MustRegister(col)
	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", c.Port), nil))
}
