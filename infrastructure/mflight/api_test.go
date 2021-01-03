package mflight_test

import (
	"fmt"
	"mflight-exporter/infrastructure/mflight"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSensorMonitor(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if path := r.URL.Path; path != "/SensorMonitorV2.xml" {
			t.Fatal(path)
		}

		if qs := r.URL.RawQuery; qs != "x-KEY_MOBILE_ID=test-mobile-id&x-KEY_UPDATE_DATE=" {
			t.Fatal(qs)
		}

		fmt.Fprint(
			w,
			`
      <db>
        <table id="67243">
          <time>202101030000</time>
          <unixtime>1609599600</unixtime>
          <temp>22.0</temp>
          <humi>43.3</humi>
          <illu>405</illu>
        </table>
        <table id="67244">
          <time>202101030005</time>
          <unixtime>1609599900</unixtime>
          <temp>21.9</temp>
          <humi>43.0</humi>
          <illu>406</illu>
        </table>
      </db>`,
		)
	}))
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
