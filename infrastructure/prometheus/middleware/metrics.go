package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	serverSv *prometheus.SummaryVec
	clientSv *prometheus.SummaryVec
)

func init() {
	serverSv = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_server_requests_seconds",
		},
		[]string{"method", "uri", "status"},
	)
	clientSv = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_client_requests_seconds",
		},
		[]string{"method", "host", "uri", "status"},
	)

	prometheus.MustRegister(serverSv, clientSv)
}
