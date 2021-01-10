package middleware

import (
	"log"
	"net/http"
	"time"
)

type HandlerMetricsMiddleware interface {
	http.Handler
}

func NewHandlerMetricsMiddleware(h http.Handler) HandlerMetricsMiddleware {
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

	log.Printf("request: %v\n", r.URL.Path)

	wr := newStatusRecoder(w)

	m.h.ServeHTTP(wr, r)

	log.Printf("response: %v\n", wr.status)

	elappsed := time.Since(start)

	s := m.s.Get(r.Method, r.URL.Path, wr.status)
	s.Observe(float64(elappsed))
}
