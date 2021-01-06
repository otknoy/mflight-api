package collector

import (
	"log"
	"mflight-api/application"

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

// MfLightCollector is the MfLight metrics collector
type MfLightCollector interface {
	prometheus.Collector
}

// NewMfLightCollector create a new MfLightCollector based on the provided MetricsCollector
func NewMfLightCollector(c application.MetricsCollector) prometheus.Collector {
	return &collector{c}
}

type collector struct {
	metricsCollector application.MetricsCollector
}

// Describe implements Collector
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- temperatureGauge.Desc()
	ch <- humidityGauge.Desc()
	ch <- illuminanceGauge.Desc()
}

// Collect implements Collector
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	m, err := c.metricsCollector.CollectMetrics()
	if err != nil {
		log.Println(err)
		return
	}

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
