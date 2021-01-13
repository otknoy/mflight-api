package mflight_test

import (
	"context"
	"errors"
	"mflight-api/domain"
	"mflight-api/infrastructure/mflight"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type stubClient struct {
	stubGetSensorMonitor func(context.Context) (*mflight.Response, error)
}

func (c *stubClient) GetSensorMonitor(ctx context.Context) (*mflight.Response, error) {
	return c.stubGetSensorMonitor(ctx)
}

func TestGetMetrics(t *testing.T) {
	c := &stubClient{
		func(ctx context.Context) (*mflight.Response, error) {
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

	m, err := sensor.GetMetrics(context.Background())
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

func TestGetMetrics_when_empty_response(t *testing.T) {
	c := &stubClient{
		func(context.Context) (*mflight.Response, error) {
			return &mflight.Response{}, nil
		},
	}

	sensor := mflight.NewMfLightSensor(c)

	m, err := sensor.GetMetrics(context.Background())
	if err == nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(domain.Metrics{}, m); diff != "" {
		t.Errorf("returned metrics is not empty\n%s", diff)
	}
}

func TestGetMetrics_when_request_failure(t *testing.T) {
	c := &stubClient{
		func(context.Context) (*mflight.Response, error) {
			return &mflight.Response{}, errors.New("test")
		},
	}

	sensor := mflight.NewMfLightSensor(c)

	m, err := sensor.GetMetrics(context.Background())
	if err == nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(domain.Metrics{}, m); diff != "" {
		t.Errorf("returned metrics is not empty\n%s", diff)
	}
}
