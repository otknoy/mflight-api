package cache_test

import (
	"mflight-api/app/infrastructure/cache"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

var minTime = time.Unix(0, 0)
var maxTime = time.Unix(1<<63-62135596801, 999999999)

func TestCacheHit(t *testing.T) {
	c := cache.New()

	c.SetWithExpiration("test-key", "test-value", maxTime)

	v, ok := c.Get("test-key")
	if !ok {
		t.Error("cache miss")
	}

	if diff := cmp.Diff("test-value", v.(string)); diff != "" {
		t.Errorf("value differs.\n%v", diff)
	}
}

func TestCacheMiss(t *testing.T) {
	c := cache.New()

	c.SetWithExpiration("test-key", "test-value", maxTime)

	_, ok := c.Get("nothing")
	if ok {
		t.Error("cache hit")
	}
}

func TestCacheHit_Expired(t *testing.T) {
	c := cache.New()

	c.SetWithExpiration("test-key", "test-value", minTime)

	_, ok := c.Get("test-key")
	if ok {
		t.Error("cache hit")
	}
}

func BenchmarkGet(b *testing.B) {
	c := cache.New()
	c.SetWithExpiration("test-key", "test-value", time.Now().Add(time.Hour))

	var wg sync.WaitGroup

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			wg.Add(1)
			go func() {
				c.Get("test-key")
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkSet(b *testing.B) {
	c := cache.New()

	var wg sync.WaitGroup

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			wg.Add(1)
			go func() {
				c.SetWithExpiration("test-key", "test-value", time.Now())
				wg.Done()
			}()
		}
		wg.Wait()
	}
	wg.Wait()
}

func BenchmarkGetSet(b *testing.B) {
	c := cache.New()

	var wg sync.WaitGroup

	get := func() {
		c.Get("test-key")
	}

	set := func() {
		c.SetWithExpiration("test-key", "test-value", time.Now())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			wg.Add(1)
			go get()
			go set()
			wg.Done()
		}
		wg.Wait()
	}
	wg.Wait()
}
