package main

import (
	"fmt"
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
	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	server := server.GracefulShutdownServer{
		Server: func() http.Server {
			c := application.NewMetricsCollector(
				mflight.NewMfLightSensor(
					mflight.NewCacheClient(
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
				),
			)

			h := handler.NewSensorMetricsHandler(c)

			col := collector.NewMfLightCollector(c)
			prometheus.MustRegister(col)

			mux := http.NewServeMux()
			mux.Handle("/getSensorMetrics", middleware.NewHandlerMetricsMiddleware(h))
			mux.Handle("/metrics", middleware.NewHandlerMetricsMiddleware(promhttp.Handler()))

			return http.Server{
				Addr:    fmt.Sprintf(":%d", config.Port),
				Handler: mux,
			}
		}(),
	}

	log.Println("server start")
	defer log.Println("server shutdown")

	if err := server.ListenAndServeWithGracefulShutdown(); err != http.ErrServerClosed {
		log.Fatal("server error: ", err)
	}
}
