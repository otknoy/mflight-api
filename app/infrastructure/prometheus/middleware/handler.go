package middleware

import (
	"net/http"
	"strconv"
	"time"
)

// InstrumentHandlerMetrics returns a middleware that Wraps the provided http.Handler
// to observe the request count and total request duration.
func InstrumentHandlerMetrics(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wr := newStatusRecoder(w)

		h.ServeHTTP(wr, r)

		duration := time.Since(start)

		serverSv.WithLabelValues(
			r.Method,
			r.URL.Path,
			strconv.Itoa(wr.status),
		).Observe(duration.Seconds())
	})
}
