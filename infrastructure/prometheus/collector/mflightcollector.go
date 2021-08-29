package collector

import (
	"context"
	"log"
	"mflight-api/domain"

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
	ctx := context.Background()

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

	temperatureGauge.Set(float64(m.Temperature))
	humidityGauge.Set(float64(m.Humidity))
	illuminanceGauge.Set(float64(m.Illuminance))

	temperatureGauge.Collect(ch)
	humidityGauge.Collect(ch)
	illuminanceGauge.Collect(ch)
}
