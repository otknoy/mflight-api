package httpclient_test

import (
	"mflight-api/infrastructure/cache"
	"mflight-api/infrastructure/mflight/httpclient"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type mockRoundTripper struct {
	http.RoundTripper
	MockRoundTrip func(*http.Request) (*http.Response, error)
}

func (rt *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt.MockRoundTrip(req)
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

func TestRoundTrip(t *testing.T) {
	u, _ := url.Parse("http://foobar.test/foo")
	testReq := &http.Request{Method: http.MethodGet, URL: u}
	res := &http.Response{}

	mockRoundTripper := &mockRoundTripper{
		MockRoundTrip: func(req *http.Request) (*http.Response, error) {
			if req != testReq {
				t.Fail()
			}
			return res, nil
		},
	}
	mockCache := &mockCache{}

	rtc := httpclient.NewRoundTripperCache(mockRoundTripper, mockCache)

	t.Run("cache miss", func(t *testing.T) {
		mockCache.MockGet = func(key string) (interface{}, bool) {
			if key != "http://foobar.test/foo" {
				t.Fail()
			}
			return nil, false
		}
		mockCache.MockSetWithExpiration = func(key string, value interface{}, _ time.Time) {
			if key != "http://foobar.test/foo" && value != res {
				t.Fail()
			}
		}

		got, err := rtc.RoundTrip(testReq)

		if err != nil {
			t.Errorf("RoundTrip returns error.\n%v", err)
		}

		if diff := cmp.Diff(res, got); diff != "" {
			t.Errorf("response differs.\n%v", diff)
		}
	})

	t.Run("cache hit", func(t *testing.T) {
		mockCache.MockGet = func(key string) (interface{}, bool) {
			if key != "http://foobar.test/foo" {
				t.Fail()
			}
			return res, true
		}
		mockCache.MockSetWithExpiration = func(_ string, _ interface{}, _ time.Time) {
			t.Fail()
		}

		got, err := rtc.RoundTrip(testReq)

		if err != nil {
			t.Errorf("RoundTrip returns error.\n%v", err)
		}

		if diff := cmp.Diff(res, got); diff != "" {
			t.Errorf("response differs.\n%v", diff)
		}
	})
}
