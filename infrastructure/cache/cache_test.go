package cache_test

import (
	"mflight-api/infrastructure/cache"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

var minTime = time.Unix(0, 0)
var maxTime = time.Unix(1<<63-62135596801, 999999999)

func TestCacheHit(t *testing.T) {
	c := cache.New()

	c.SetWithExpiration("test-key", "test-value", maxTime)

	v := c.Get("test-key")

	if diff := cmp.Diff("test-value", v.(string)); diff != "" {
		t.Errorf("value differs.\n%v", diff)
	}
}

func TestCacheMiss(t *testing.T) {
	c := cache.New()

	c.SetWithExpiration("test-key", "test-value", maxTime)

	v := c.Get("nothing")

	if v != nil {
		t.Errorf("value should be nil. but %v", v)
	}
}

func TestCacheHit_Expired(t *testing.T) {
	c := cache.New()

	c.SetWithExpiration("test-key", "test-value", minTime)

	v := c.Get("test-key")

	if v != nil {
		t.Errorf("value should be nil. but %v", v)
	}
}
