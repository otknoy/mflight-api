package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpReq = prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "http_requests_seconds",
	})
)

func init() {
	prometheus.MustRegister(httpReq)
}

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
	start := time.Now()

	log.Printf("request: %v\n", r.URL.Path)

	wr := &statusRecoder{
		status: 0,
		w:      w,
	}

	m.h.ServeHTTP(wr, r)

	log.Printf("response: %v\n", wr.status)

	elappsed := time.Since(start)
	httpReq.Observe(float64(elappsed))
}

type statusRecoder struct {
	status int
	w      http.ResponseWriter
}

func (r *statusRecoder) Header() http.Header {
	return r.w.Header()
}

func (r *statusRecoder) Write(bytes []byte) (int, error) {
	return r.w.Write(bytes)
}

func (r *statusRecoder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.w.WriteHeader(statusCode)
}
