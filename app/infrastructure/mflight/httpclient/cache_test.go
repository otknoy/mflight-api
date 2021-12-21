package httpclient_test

import (
	"context"
	"mflight-api/app/infrastructure/cache"
	"mflight-api/app/infrastructure/mflight/httpclient"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type mockClient struct {
	httpclient.Client
	MockGetSensorMonitor func(ctx context.Context) (*httpclient.Response, error)
}

func (c *mockClient) GetSensorMonitor(ctx context.Context) (*httpclient.Response, error) {
	return c.MockGetSensorMonitor(ctx)
}

type mockCache struct {
	cache.Cache
	MockGet               func(key string) (interface{}, bool)
	MockSetWithExpiration func(key string, value interface{}, expiration time.Time)
}

func (c *mockCache) Get(key string) (interface{}, bool) {
	return c.MockGet(key)
}

func (c *mockCache) SetWithExpiration(key string, value interface{}, expiration time.Time) {
	c.MockSetWithExpiration(key, value, expiration)
}

func TestCacheClient_GetSensorMonitor(t *testing.T) {
	testCtx := context.Background()

	res := &httpclient.Response{
		Tables: []httpclient.Table{
			{
				Temperature: 23.4,
				Humidity:    45.6,
				Illuminance: 678,
			},
		},
	}

	mockClient := &mockClient{
		MockGetSensorMonitor: func(ctx context.Context) (*httpclient.Response, error) {
			if ctx != testCtx {
				t.Fail()
			}
			return res, nil
		},
	}
	mockCache := &mockCache{}

	c := httpclient.NewCacheClient(mockClient, mockCache, 5*time.Second)

	t.Run("cache miss", func(t *testing.T) {
		mockCache.MockGet = func(key string) (interface{}, bool) {
			if key != "fixed" {
				t.Fail()
			}
			return nil, false
		}
		mockCache.MockSetWithExpiration = func(key string, value interface{}, _ time.Time) {
			if key != "fixed" && value != res {
				t.Fail()
			}
		}

		got, err := c.GetSensorMonitor(testCtx)

		if err != nil {
			t.Errorf("GetSensorMonitor returns error.\n%v", err)
		}

		if diff := cmp.Diff(res, got); diff != "" {
			t.Errorf("response differs.\n%v", diff)
		}
	})

	t.Run("cache hit", func(t *testing.T) {
		mockCache.MockGet = func(key string) (interface{}, bool) {
			if key != "fixed" {
				t.Fail()
			}
			return res, true
		}

		got, err := c.GetSensorMonitor(testCtx)

		if err != nil {
			t.Errorf("GetSensorMonitor returns error.\n%v", err)
		}

		if diff := cmp.Diff(res, got); diff != "" {
			t.Errorf("response differs.\n%v", diff)
		}
	})
}
