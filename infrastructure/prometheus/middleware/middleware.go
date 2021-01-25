package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	serverSv *prometheus.SummaryVec
)

func init() {
	serverSv = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_server_requests_seconds",
		},
		[]string{"method", "uri", "status"},
	)

	prometheus.MustRegister(serverSv)
}
