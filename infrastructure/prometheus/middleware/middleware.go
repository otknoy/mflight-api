package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	sv *prometheus.SummaryVec
)

func init() {
	sv = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_requests_seconds",
		},
		[]string{"method", "uri", "status"},
	)

	prometheus.MustRegister(sv)
}
