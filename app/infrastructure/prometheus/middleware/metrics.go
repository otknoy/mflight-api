package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	serverSv *prometheus.SummaryVec
	clientSv *prometheus.SummaryVec
)

func init() {
	serverSv = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_server_requests_seconds",
		},
		[]string{"method", "uri", "status"},
	)
	clientSv = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_client_requests_seconds",
		},
		[]string{"method", "host", "uri", "status"},
	)
}
