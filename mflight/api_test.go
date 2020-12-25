package mflight_test

import (
	"mflight-exporter/mflight"
	"testing"
)

func TestBuildURL(t *testing.T) {
	url := mflight.BuildURL("http://example.com:8080", "test-mobile-id")

	want := "http://example.com:8080/SensorMonitorV2.xml?x-KEY_MOBILE_ID=test-mobile-id&x-KEY_UPDATE_DATE="

	if url != want {
		t.Errorf("\n%v\n%v\n", url, want)
	}
}
