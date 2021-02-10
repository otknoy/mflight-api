package main

import (
	"log"
	"mflight-api/application"
	"mflight-api/config"
	"mflight-api/infrastructure/cache"
	"mflight-api/infrastructure/mflight"
	"mflight-api/infrastructure/prometheus/collector"
	"mflight-api/infrastructure/prometheus/middleware"
	"mflight-api/interfaces/handler"
	"mflight-api/interfaces/server"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	c, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	metricsCollector := initMetricsCollector(&c.MfLight)

	h := handler.NewSensorMetricsHandler(metricsCollector)

	col := collector.NewMfLightCollector(metricsCollector)
	prometheus.MustRegister(col)

	s := server.NewServer(
		map[string]http.Handler{
			"/getSensorMetrics": middleware.NewHandlerMetricsMiddleware(h),
			"/metrics":          middleware.NewHandlerMetricsMiddleware(promhttp.Handler()),
		},
		c.Port,
	)

	idleConnsClosed := make(chan struct{})
	s.ListenAndServe(idleConnsClosed)
	<-idleConnsClosed
}

func initMetricsCollector(c *config.MfLightConfig) application.MetricsCollector {
	sensor := mflight.NewMfLightSensor(
		mflight.NewCacheClient(
			mflight.NewClient(
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
