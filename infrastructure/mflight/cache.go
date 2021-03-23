package mflight

import (
	"context"
	"time"

	"mflight-api/infrastructure/cache"
	"mflight-api/infrastructure/mflight/httpclient"
)

type ClientFunc func(context.Context) (*httpclient.Response, error)

func (f ClientFunc) GetSensorMonitor(ctx context.Context) (*httpclient.Response, error) {
	return f(ctx)
}

// NewCacheClient wraps client to enable caching
func NewCacheClient(client httpclient.Client, cache cache.Cache, ttl time.Duration) httpclient.Client {
	const key = "fixed"

	return ClientFunc(func(ctx context.Context) (*httpclient.Response, error) {
		v := cache.Get(key)
		if v != nil {
			return v.(*httpclient.Response), nil
		}

		r, err := client.GetSensorMonitor(ctx)

		if err == nil {
			cache.SetWithExpiration(key, r, time.Now().Add(ttl))
		}

		return r, err
	})
}
