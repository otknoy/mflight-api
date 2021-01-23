package main

import (
	"context"
	"fmt"
	"log"
	"mflight-api/application"
	"mflight-api/config"
	"mflight-api/infrastructure/mflight"
	"mflight-api/infrastructure/prometheus/collector"
	"mflight-api/infrastructure/prometheus/middleware"
	"mflight-api/interfaces/handler"
	"net/http"
	"os"
	"os/signal"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	c, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	sensor := mflight.NewMfLightSensor(
		mflight.NewCacheClient(
			mflight.NewClient(c.MfLight.URL, c.MfLight.MobileID),
		),
	)
	metricsCollector := application.NewMetricsCollector(sensor)

	h := handler.NewSensorMetricsHandler(metricsCollector)

	col := collector.NewMfLightCollector(metricsCollector)
	prometheus.MustRegister(col)

	mux := http.NewServeMux()
	mux.Handle("/getSensorMetrics", middleware.NewHandlerMetricsMiddleware(h))
	mux.Handle("/metrics", middleware.NewHandlerMetricsMiddleware(promhttp.Handler()))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", c.Port),
		Handler: mux,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}

		close(idleConnsClosed)
	}()
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	<-idleConnsClosed
}
