package mflight_test

import (
	"context"
	"mflight-api/infrastructure/mflight"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetSensorMonitor(t *testing.T) {
	s := NewStubServer(t, `
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
      </db>`)
	defer s.Close()

	c := mflight.NewClient(s.URL, "test-mobile-id")
	res, err := c.GetSensorMonitor(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	want := &mflight.Response{
		Tables: []mflight.Table{
			{
				ID:          67243,
				Time:        "202101030000",
				Unixtime:    1609599600,
				Temperature: 22.0,
				Humidity:    43.3,
				Illuminance: 405,
			},
			{
				ID:          67244,
				Time:        "202101030005",
				Unixtime:    1609599900,
				Temperature: 21.9,
				Humidity:    43.0,
				Illuminance: 406,
			},
		},
	}
	if diff := cmp.Diff(want, res); diff != "" {
		t.Errorf("response differs.\n%v", diff)
	}
}

func TestGetSensorMonitor_network_error(t *testing.T) {
	c := mflight.NewClient("http://hoge.test", "dummy-mobile-id")

	res, err := c.GetSensorMonitor(context.Background())

	if err != nil {
		t.Errorf("should return error.")
	}
	if len(res.Tables) != 0 {
		t.Errorf("shoud return empty response.")
	}
}

func TestGetSensorMonitor_invalid_response(t *testing.T) {
	s := NewStubServer(t, `<db><table</table></db>`)
	defer s.Close()

	c := mflight.NewClient(s.URL, "test-mobile-id")
	res, err := c.GetSensorMonitor(context.Background())

	if err != nil {
		t.Errorf("should return error.")
	}
	if len(res.Tables) != 0 {
		t.Errorf("shoud return empty response.")
	}
}

func TestBuildRequestWithContext(t *testing.T) {
	r := mflight.BuildRequestWithContext(context.Background(), "http://example.com:8080", "test-mobile-id")

	want := "http://example.com:8080/SensorMonitorV2.xml?x-KEY_MOBILE_ID=test-mobile-id&x-KEY_UPDATE_DATE="

	if diff := cmp.Diff(want, r.URL.String()); diff != "" {
		t.Errorf("request differs\n%s", diff)
	}
}
