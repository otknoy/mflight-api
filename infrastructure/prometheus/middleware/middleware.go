package middleware

import (
	"net/http"
	"strconv"
	"time"

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

// NewHandlerMetricsMiddleware returns a middleware that Wraps the provided http.Handler
// to observe the request count and total request duration.
func NewHandlerMetricsMiddleware(h http.Handler) http.Handler {
	return &middleware{
		h: h,
	}
}

type middleware struct {
	h http.Handler
}

func (m *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	wr := newStatusRecoder(w)

	m.h.ServeHTTP(wr, r)

	duration := time.Since(start)

	sv.WithLabelValues(
		r.Method,
		r.URL.Path,
		strconv.Itoa(wr.status),
	).Observe(duration.Seconds())
}
