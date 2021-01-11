package middleware

import (
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type summaries map[string]prometheus.Summary

func (s summaries) Get(method, uri string, status int) prometheus.Summary {
	k := fmt.Sprintf("%s%s%d", method, uri, status)

	v, ok := s[k]
	if ok {
		return v
	}

	s[k] = prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "http_requests_seconds",
		ConstLabels: prometheus.Labels{
			"method": method,
			"uri":    uri,
			"status": strconv.Itoa(status),
		},
	})
	prometheus.MustRegister(s[k])

	return s[k]
}
