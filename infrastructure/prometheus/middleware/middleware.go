package middleware

import (
	"net/http"
	"time"
)

// NewHandlerMetricsMiddleware returns a middleware that Wraps the provided http.Handler
// to observe the request count and total request duration.
func NewHandlerMetricsMiddleware(h http.Handler) http.Handler {
	return &middleware{
		s: make(summaries),
		h: h,
	}
}

type middleware struct {
	s summaries
	h http.Handler
}

func (m *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	wr := newStatusRecoder(w)

	m.h.ServeHTTP(wr, r)

	elappsed := time.Since(start)

	s := m.s.Get(r.Method, r.URL.Path, wr.status)
	s.Observe(float64(elappsed))
}
