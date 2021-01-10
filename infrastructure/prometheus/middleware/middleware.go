package middleware

import (
	"log"
	"net/http"
)

type HandlerMetricsMiddleware interface {
	http.Handler
}

func NewHandlerMetricsMiddleware(h http.Handler) HandlerMetricsMiddleware {
	return &middleware{h}
}

type middleware struct {
	h http.Handler
}

func (m *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("request: %v\n", r.URL.Path)

	m.h.ServeHTTP(w, r)

	log.Printf("response: %v\n", w)
}
