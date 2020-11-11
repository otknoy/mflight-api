package collector

import (
	"mflight-exporter/mflight"

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

type MfLightCollector struct {
	sensor mflight.Sensor
}

func NewMfLightCollector(sensor mflight.Sensor) prometheus.Collector {
	return &MfLightCollector{sensor}
}

func (c *MfLightCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- temperatureGauge.Desc()
	ch <- humidityGauge.Desc()
	ch <- illuminanceGauge.Desc()
}

func (c *MfLightCollector) Collect(ch chan<- prometheus.Metric) {
	m, _ := c.sensor.GetMetrics()

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
