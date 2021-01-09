package mflight_test

import (
	"mflight-api/domain"
	"mflight-api/infrastructure/mflight"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type stubClient struct {
	stubGetSensorMonitor func() (*mflight.Response, error)
}

func (c *stubClient) GetSensorMonitor() (*mflight.Response, error) {
	return c.stubGetSensorMonitor()
}

func TestGetMetrics(t *testing.T) {
	c := &stubClient{
		func() (*mflight.Response, error) {
			return &mflight.Response{
				Tables: []mflight.Table{
					{
						Temperature: 25.4,
						Humidity:    65.7,
						Illuminance: 234,
					},
					{
						Temperature: 21.9,
						Humidity:    43.0,
						Illuminance: 406,
					},
				},
			}, nil
		},
	}

	sensor := mflight.NewMfLightSensor(c)

	m, err := sensor.GetMetrics()
	if err != nil {
		t.Fatal(err)
	}

	want := domain.Metrics{
		Temperature: 21.9,
		Humidity:    43.0,
		Illuminance: 406,
	}
	if diff := cmp.Diff(want, m); diff != "" {
		t.Errorf("returned metrics differs\n%s", diff)
	}
}
