package middleware

import (
	"net/http"
	"time"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

// InstrumentRoundTripperMetrics returns a middleware that Wraps the provided http.RoundTripper
// to observe the request count and total request duration.
func InstrumentRoundTripperMetrics(rt http.RoundTripper) http.RoundTripper {
	return roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		start := time.Now()

		res, err := rt.RoundTrip(r)

		duration := time.Since(start)

		clientSv.WithLabelValues(
			r.Method,
			r.URL.Host,
			r.URL.Path,
			status(res),
		).Observe(duration.Seconds())

		return res, err
	})
}

func status(r *http.Response) string {
	if r == nil {
		return "CLIENT_ERROR"
	}
	return r.Status
}
