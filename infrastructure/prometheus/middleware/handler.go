package middleware

import (
	"net/http"
	"strconv"
	"time"
)

// NewHandlerMetricsMiddleware returns a middleware that Wraps the provided http.Handler
// to observe the request count and total request duration.
func NewHandlerMetricsMiddleware(h http.Handler) http.Handler {
	return &httpHandlerMiddleware{
		h: h,
	}
}

type httpHandlerMiddleware struct {
	h http.Handler
}

var _ http.Handler = (*httpHandlerMiddleware)(nil)

func (m *httpHandlerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	wr := newStatusRecoder(w)

	m.h.ServeHTTP(wr, r)

	duration := time.Since(start)

	serverSv.WithLabelValues(
		r.Method,
		r.URL.Path,
		strconv.Itoa(wr.status),
	).Observe(duration.Seconds())
}
