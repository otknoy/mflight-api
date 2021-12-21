package config_test

import (
	"mflight-api/app/config"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestLoad(t *testing.T) {
	t.Setenv("APP_PORT", "5000")
	t.Setenv("APP_MFLIGHT_URL", "http://example.com:56002")
	t.Setenv("APP_MFLIGHT_MOBILE_ID", "test-mobileID")
	t.Setenv("APP_MFLIGHT_CACHE_TTL", "15s")

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
		t.Setenv("APP_PORT", "5000")
		t.Setenv("APP_MFLIGHT_URL", "http://example.com:56002")
		t.Setenv("APP_MFLIGHT_MOBILE_ID", "test-mobileID")
		t.Setenv("APP_MFLIGHT_CACHE_TTL", tt.ttl)

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

func TestLoad_return_error_if_port_number_is_invalid(t *testing.T) {
	t.Setenv("APP_PORT", "invalid-port-number")

	_, err := config.Load()
	if err == nil {
		t.Error("error should occur.")
	}
}
