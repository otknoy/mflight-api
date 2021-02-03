package config_test

import (
	"mflight-api/config"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestLoad(t *testing.T) {
	os.Setenv("APP_PORT", "5000")
	os.Setenv("APP_MFLIGHT_URL", "http://example.com:56002")
	os.Setenv("APP_MFLIGHT_MOBILE_ID", "test-mobileID")
	os.Setenv("APP_MFLIGHT_CACHE_TTL", "15s")

	got, _ := config.Load()

	want := config.AppConfig{
		Port: 5000,
		MfLight: config.MfLightConfig{
			URL:      "http://example.com:56002",
			MobileID: "test-mobileID",
			CacheTTL: 15 * time.Second,
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("config differs.\n%v", diff)
	}
}

func TestLoad_cache_ttl(t *testing.T) {
	tests := []struct {
		ttl  string
		want time.Duration
	}{
		{"1ms", time.Millisecond},
		{"2s", 2 * time.Second},
		{"3m", 3 * time.Minute},
		{"", 0},
		{"invalid-duration-string", 0},
	}

	for _, tt := range tests {
		os.Setenv("APP_PORT", "5000")
		os.Setenv("APP_MFLIGHT_URL", "http://example.com:56002")
		os.Setenv("APP_MFLIGHT_MOBILE_ID", "test-mobileID")
		os.Setenv("APP_MFLIGHT_CACHE_TTL", tt.ttl)

		got, _ := config.Load()

		want := config.AppConfig{
			Port: 5000,
			MfLight: config.MfLightConfig{
				URL:      "http://example.com:56002",
				MobileID: "test-mobileID",
				CacheTTL: tt.want,
			},
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("config differs.\n%v", diff)
		}
	}
}
