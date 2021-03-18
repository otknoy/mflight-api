package mflight_test

import (
	"context"
	"mflight-api/infrastructure/cache"
	"mflight-api/infrastructure/mflight"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type mockClient struct {
	mflight.Client
	MockGetSensorMonitor func(ctx context.Context) (*mflight.Response, error)
}

func (c *mockClient) GetSensorMonitor(ctx context.Context) (*mflight.Response, error) {
	return c.MockGetSensorMonitor(ctx)
}

type mockCache struct {
	cache.Cache
	MockGet               func(key string) interface{}
	MockSetWithExpiration func(key string, value interface{}, expiration time.Time)
}

func (c *mockCache) Get(key string) interface{} {
	return c.MockGet(key)
}

func (c *mockCache) SetWithExpiration(key string, value interface{}, expiration time.Time) {
	c.MockSetWithExpiration(key, value, expiration)
}

func TestCacheClient_GetSensorMonitor(t *testing.T) {
	testCtx := context.Background()

	res := &mflight.Response{
		Tables: []mflight.Table{
			{
				Temperature: 23.4,
				Humidity:    45.6,
				Illuminance: 678,
			},
		},
	}

	mockClient := &mockClient{
		MockGetSensorMonitor: func(ctx context.Context) (*mflight.Response, error) {
			if ctx != testCtx {
				t.Fail()
			}
			return res, nil
		},
	}
	mockCache := &mockCache{}

	c := mflight.NewCacheClient(mockClient, mockCache, 5*time.Second)

	t.Run("cache miss", func(t *testing.T) {
		mockCache.MockGet = func(key string) interface{} {
			if key != "fixed" {
				t.Fail()
			}
			return nil
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
		mockCache.MockGet = func(key string) interface{} {
			if key != "fixed" {
				t.Fail()
			}
			return res
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
