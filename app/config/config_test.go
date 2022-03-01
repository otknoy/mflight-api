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

	got, err := config.Load()

	if err != nil {
		t.Errorf("error should not occur.\n%v", err)
	}

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

func TestLoad_returns_error_when_invalid_configuration(t *testing.T) {
	tests := []struct {
		port, cacheTtl string
	}{
		{"5000", "invalid-cache-ttl"},
		{"invalid-port-number", "15s"},
	}

	for _, tt := range tests {
		t.Setenv("APP_PORT", tt.port)
		t.Setenv("APP_MFLIGHT_URL", "http://example.com:56002")
		t.Setenv("APP_MFLIGHT_MOBILE_ID", "test-mobileID")
		t.Setenv("APP_MFLIGHT_CACHE_TTL", tt.cacheTtl)

		_, err := config.Load()

		if err == nil {
			t.Errorf("error should occur.\n%v", tt)
		}
	}
}
