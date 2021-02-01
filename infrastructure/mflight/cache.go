package mflight

import (
	"context"
	"time"

	"mflight-api/infrastructure/cache"
)

type ClientFunc func(context.Context) (*Response, error)

func (f ClientFunc) GetSensorMonitor(ctx context.Context) (*Response, error) {
	return f(ctx)
}

// NewCacheClient wraps client to enable caching
func NewCacheClient(client Client, cache cache.Cache) Client {
	const key = "fixed"

	return ClientFunc(func(ctx context.Context) (*Response, error) {
		v := cache.Get(key)
		if v != nil {
			return v.(*Response), nil
		}

		r, err := client.GetSensorMonitor(ctx)

		if err == nil {
			cache.SetWithExpiration(key, r, time.Now().Add(5*time.Second))
		}

		return r, err
	})
}
