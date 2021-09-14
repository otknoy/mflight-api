package httpclient

import (
	"mflight-api/infrastructure/cache"
	"net/http"
	"time"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func NewRoundTripperCache(rt http.RoundTripper, cache cache.Cache) http.RoundTripper {
	ttl := 5 * time.Second

	return roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodGet {
			return rt.RoundTrip(r)
		}

		key := r.URL.String()

		v, ok := cache.Get(key)
		if ok {
			return v.(*http.Response), nil
		}

		res, err := rt.RoundTrip(r)
		if err == nil {
			cache.SetWithExpiration(key, res, time.Now().Add(ttl))
		}

		return res, err
	})
}
