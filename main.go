package main

import (
	"context"
	"log"
	"mflight-api/application"
	"mflight-api/config"
	"mflight-api/infrastructure/cache"
	"mflight-api/infrastructure/mflight"
	"mflight-api/infrastructure/mflight/httpclient"
	"mflight-api/infrastructure/prometheus/collector"
	"mflight-api/infrastructure/prometheus/middleware"
	"mflight-api/infrastructure/server"
	"mflight-api/interfaces/handler"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	c, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	mc := initMetricsCollector(&c.MfLight)

	s := initServer(c.Port, mc)

	log.Println("server start")

	if err := s.ListenAndServeWithGracefulShutdown(context.Background()); err != nil {
		log.Fatalln("HTTP server Shutdown: ", err)
	}

	log.Println("server shutdown")
}

func initServer(port int, mc application.MetricsCollector) *server.GracefulShutdownServer {
	h := handler.NewSensorMetricsHandler(mc)

	col := collector.NewMfLightCollector(mc)
	prometheus.MustRegister(col)

	mux := http.NewServeMux()
	mux.Handle("/getSensorMetrics", middleware.NewHandlerMetricsMiddleware(h))
	mux.Handle("/metrics", middleware.NewHandlerMetricsMiddleware(promhttp.Handler()))

	return server.NewServer(mux, port)
}

func initMetricsCollector(c *config.MfLightConfig) application.MetricsCollector {
	sensor := mflight.NewMfLightSensor(
		mflight.NewCacheClient(
			httpclient.NewClient(
				&http.Client{
					Transport: middleware.NewRoundTripperMetricsMiddleware(http.DefaultTransport),
				},
				c.URL,
				c.MobileID,
			),
			cache.New(),
			c.CacheTTL,
		),
	)

	return application.NewMetricsCollector(sensor)
}
