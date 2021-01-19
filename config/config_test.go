package config_test

import (
	"mflight-api/config"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoad(t *testing.T) {
	os.Setenv("APP_PORT", "5000")
	os.Setenv("APP_MFLIGHT_URL", "http://example.com:56002")
	os.Setenv("APP_MFLIGHT_MOBILE_ID", "test-mobileID")

	got, _ := config.Load()

	want := config.AppConfig{
		Port: 5000,
		MfLight: config.MfLightConfig{
			URL:      "http://example.com:56002",
			MobileID: "test-mobileID",
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("config differs.\n%v", diff)
	}
}
