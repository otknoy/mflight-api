package main

import (
	"fmt"
	"mflight-api/app/config"
	"mflight-api/app/handler"
	"mflight-api/app/infrastructure/cache"
	"mflight-api/app/infrastructure/mflight"
	"mflight-api/app/infrastructure/mflight/httpclient"
	"mflight-api/app/infrastructure/prometheus/collector"
	"mflight-api/app/infrastructure/prometheus/middleware"
	"mflight-api/app/infrastructure/server"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func init() {
	l, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(l)
}

func main() {
	config, err := config.Load()
	if err != nil {
		zap.L().Fatal("config load failure", zap.Error(err))
	}

	server := initServer(config)

	zap.L().Info("server start")
	defer zap.L().Info("server shutdown")

	if err := server.ListenAndServeGracefully(); err != http.ErrServerClosed {
		zap.L().Fatal("server error", zap.Error(err))
	}
}

func initServer(config config.AppConfig) server.GracefulShutdownServer {
	metricsGetter := mflight.NewMetricsGetter(
		httpclient.NewCacheClient(
			httpclient.NewClient(
				&http.Client{
					Transport: middleware.InstrumentRoundTripperMetrics(http.DefaultTransport),
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
	mux.Handle("/getSensorMetrics", middleware.InstrumentHandlerMetrics(h))
	mux.Handle("/metrics", middleware.InstrumentHandlerMetrics(promhttp.Handler()))

	return server.GracefulShutdownServer{
		Server: http.Server{
			Addr:    fmt.Sprintf(":%d", config.Port),
			Handler: mux,
		},
	}
}
