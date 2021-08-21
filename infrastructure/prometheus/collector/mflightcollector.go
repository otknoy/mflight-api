package collector

import (
	"context"
	"log"
	"mflight-api/domain"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "MultiFunctionLight"
)

var (
	temperatureGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "temperature",
		Help:      "multifunction light temperature",
	})
	humidityGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "humidity",
		Help:      "multifunction light humidity",
	})
	illuminanceGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "illuminance",
		Help:      "multifunction light illuminance",
	})
)

// NewMfLightCollector create a new prometheus.Collector based on the provided MetricsGetter
func NewMfLightCollector(g domain.MetricsGetter) prometheus.Collector {
	return &collector{g}
}

type collector struct {
	metricsGetter domain.MetricsGetter
}

// Describe implements Collector
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- temperatureGauge.Desc()
	ch <- humidityGauge.Desc()
	ch <- illuminanceGauge.Desc()
}

// Collect implements Collector
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	mch := make(chan domain.Metrics)
	defer close(mch)

	go func() {
		l, err := c.metricsGetter.GetMetrics(ctx)
		if err != nil {
			log.Printf("failed to collect metrics: %v", err)
			return
		}

		m, err := l.Last()
		if err != nil {
			log.Printf("failed to collect metrics: %v", err)
			return
		}

		mch <- m
	}()

	select {
	case <-ctx.Done():
		log.Println("timeout: ", ctx.Err())
	case m := <-mch:
		ch <- prometheus.MustNewConstMetric(
			temperatureGauge.Desc(),
			prometheus.GaugeValue,
			float64(m.Temperature),
		)
		ch <- prometheus.MustNewConstMetric(
			humidityGauge.Desc(),
			prometheus.GaugeValue,
			float64(m.Humidity),
		)
		ch <- prometheus.MustNewConstMetric(
			illuminanceGauge.Desc(),
			prometheus.GaugeValue,
			float64(m.Illuminance),
		)
	}
}
