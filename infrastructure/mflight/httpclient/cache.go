package httpclient

import (
	"context"
	"time"

	"mflight-api/infrastructure/cache"
)

type clientFunc func(context.Context) (*Response, error)

func (f clientFunc) GetSensorMonitor(ctx context.Context) (*Response, error) {
	return f(ctx)
}

// NewCacheClient wraps client to enable caching
func NewCacheClient(client Client, cache cache.Cache, ttl time.Duration) Client {
	const key = "fixed"

	return clientFunc(func(ctx context.Context) (*Response, error) {
		v, ok := cache.Get(key)
		if ok {
			return v.(*Response), nil
		}

		r, err := client.GetSensorMonitor(ctx)
		if err == nil {
			cache.SetWithExpiration(key, r, time.Now().Add(ttl))
		}

		return r, err
	})
}
