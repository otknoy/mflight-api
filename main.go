package main

import (
	"context"
	"fmt"
	"log"
	"mflight-api/config"
	"mflight-api/handler"
	"mflight-api/infrastructure/cache"
	"mflight-api/infrastructure/mflight"
	"mflight-api/infrastructure/mflight/httpclient"
	"mflight-api/infrastructure/prometheus/collector"
	"mflight-api/infrastructure/prometheus/middleware"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	server := func() http.Server {
		metricsGetter := mflight.NewMfLightSensor(
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
		)

		h := handler.NewSensorMetricsHandler(metricsGetter)

		prometheus.MustRegister(
			collector.NewMfLightCollector(metricsGetter),
		)

		mux := http.NewServeMux()
		mux.Handle("/getSensorMetrics", middleware.NewHandlerMetricsMiddleware(h))
		mux.Handle("/metrics", middleware.NewHandlerMetricsMiddleware(promhttp.Handler()))

		return http.Server{
			Addr:    fmt.Sprintf(":%d", config.Port),
			Handler: mux,
		}
	}()

	idleConnsClosed := make(chan struct{})
	go func() {
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Println("shutdown gracefully...")
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown error: %v", err)
		}

		close(idleConnsClosed)
	}()

	log.Println("server start")
	defer log.Println("server shutdown")

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal("server error: ", err)
	}

	<-idleConnsClosed
}
