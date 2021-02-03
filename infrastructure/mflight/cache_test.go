package mflight_test

import (
	"context"
	"mflight-api/infrastructure/cache/mock_cache"
	"mflight-api/infrastructure/mflight"
	"mflight-api/infrastructure/mflight/mock_mflight"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

func TestCacheClient_GetSensorMonitor(t *testing.T) {
	ctx := context.Background()
	want := &mflight.Response{
		Tables: []mflight.Table{
			{
				Temperature: 23.4,
				Humidity:    45.6,
				Illuminance: 678,
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_mflight.NewMockClient(ctrl)
	mockCache := mock_cache.NewMockCache(ctrl)

	mockClient.EXPECT().
		GetSensorMonitor(ctx).Return(want, nil)

	c := mflight.NewCacheClient(mockClient, mockCache, 5*time.Second)

	t.Run("cache miss", func(t *testing.T) {
		mockCache.EXPECT().
			Get("fixed").Return(nil)
		mockCache.EXPECT().
			SetWithExpiration("fixed", want, gomock.Any())

		got, err := c.GetSensorMonitor(ctx)

		if err != nil {
			t.Errorf("GetSensorMonitor returns error.\n%v", err)
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("response differs.\n%v", diff)
		}

	})

	t.Run("cache hit", func(t *testing.T) {
		mockCache.EXPECT().
			Get("fixed").Return(want)

		got, err := c.GetSensorMonitor(ctx)

		if err != nil {
			t.Errorf("GetSensorMonitor returns error.\n%v", err)
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("response differs.\n%v", diff)
		}
	})
}
