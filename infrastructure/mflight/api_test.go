package mflight_test

import (
	"mflight-exporter/infrastructure/mflight"
	"testing"
)

func TestGetSensorMonitor(t *testing.T) {
	s := NewStubServer(t)
	defer s.Close()

	res, err := mflight.GetSensorMonitor(s.URL, "test-mobile-id")
	if err != nil {
		t.Fatal(err)
	}

	if len := len(res.Tables); len != 2 {
		t.Errorf("table length expect 2, but %d\n", len)
	}
	if v := res.Tables[0].Temperature; v != 22.0 {
		t.Errorf("invalid temperature: %v\n", v)
	}
	if v := res.Tables[0].Humidity; v != 43.3 {
		t.Errorf("invalid humidity: %v\n", v)
	}
	if v := res.Tables[0].Illuminance; v != 405 {
		t.Errorf("invalid illuminance: %v\n", v)
	}
}

func TestBuildURL(t *testing.T) {
	url := mflight.BuildURL("http://example.com:8080", "test-mobile-id")

	want := "http://example.com:8080/SensorMonitorV2.xml?x-KEY_MOBILE_ID=test-mobile-id&x-KEY_UPDATE_DATE="

	if url != want {
		t.Errorf("\n%v\n%v\n", url, want)
	}
}
