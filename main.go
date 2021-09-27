package main

import (
	"fmt"
	"log"
	"mflight-api/config"
	"mflight-api/handler"
	"mflight-api/infrastructure/cache"
	"mflight-api/infrastructure/mflight"
	"mflight-api/infrastructure/mflight/httpclient"
	"mflight-api/infrastructure/prometheus/collector"
	"mflight-api/infrastructure/prometheus/middleware"
	"mflight-api/infrastructure/server"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	server := initServer(config)

	log.Println("server start")
	defer log.Println("server shutdown")

	if err := server.ListenAndServeGracefully(); err != http.ErrServerClosed {
		log.Fatal("server error: ", err)
	}
}

func initServer(config config.AppConfig) server.GracefulShutdownServer {
	metricsGetter := mflight.NewMetricsGetter(
		httpclient.NewCacheClient(
			httpclient.NewClient(
				&http.Client{
					Transport: middleware.NewRoundTripperMetricsMiddleware(http.DefaultTransport),
				},
				config.MfLight.URL,
				config.MfLight.MobileID,
			),
			cache.New(),
			config.MfLight.CacheTTL,
		),
	)

	h := handler.NewSensorMetricsHandler(metricsGetter)

	prometheus.MustRegister(
		collector.NewMfLightCollector(metricsGetter),
	)

	mux := http.NewServeMux()
	mux.Handle("/getSensorMetrics", middleware.NewHandlerMetricsMiddleware(h))
	mux.Handle("/metrics", middleware.NewHandlerMetricsMiddleware(promhttp.Handler()))

	return server.GracefulShutdownServer{
		Server: http.Server{
			Addr:    fmt.Sprintf(":%d", config.Port),
			Handler: mux,
		},
	}
}
