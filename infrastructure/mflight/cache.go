package mflight

import (
	"context"
	"time"

	"mflight-api/infrastructure/cache"
)

// NewCacheClient wraps client to enable caching
func NewCacheClient(client Client, cache cache.Cache) Client {
	return &cacheClient{client, cache}
}

const key = "fixed"

type cacheClient struct {
	client Client
	cache  cache.Cache
}

func (c *cacheClient) GetSensorMonitor(ctx context.Context) (*Response, error) {
	v := c.cache.Get(key)
	if v != nil {
		return v.(*Response), nil
	}

	r, err := c.client.GetSensorMonitor(ctx)
	if err != nil {
		return r, err
	}

	c.cache.SetWithExpiration(key, r, time.Now().Add(5*time.Second))

	return r, err
}
