package middleware

import (
	"net/http"
	"time"
)

// NewRoundTripperMetricsMiddleware returns a middleware that Wraps the provided http.RoundTripper
// to observe the request count and total request duration.
func NewRoundTripperMetricsMiddleware(rt http.RoundTripper) http.RoundTripper {
	return &httpRoundTripperMiddleware{
		rt: rt,
	}
}

type httpRoundTripperMiddleware struct {
	rt http.RoundTripper
}

var _ http.RoundTripper = (*httpRoundTripperMiddleware)(nil)

func (m *httpRoundTripperMiddleware) RoundTrip(r *http.Request) (*http.Response, error) {
	start := time.Now()

	res, err := m.rt.RoundTrip(r)

	duration := time.Since(start)

	clientSv.WithLabelValues(
		r.Method,
		r.URL.Host,
		r.URL.Path,
		res.Status,
	).Observe(duration.Seconds())

	return res, err
}
