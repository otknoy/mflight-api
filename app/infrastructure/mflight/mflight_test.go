package mflight_test

import (
	"context"
	"errors"
	"mflight-api/app/domain"
	"mflight-api/app/infrastructure/mflight"
	"mflight-api/app/infrastructure/mflight/httpclient"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type stubClient struct {
	httpclient.Client
	MockGetSensorMonitor func(context.Context) (*httpclient.Response, error)
}

func (c *stubClient) GetSensorMonitor(ctx context.Context) (*httpclient.Response, error) {
	return c.MockGetSensorMonitor(ctx)
}

func TestGetMetrics(t *testing.T) {
	g := mflight.NewMetricsGetter(&stubClient{
		MockGetSensorMonitor: func(ctx context.Context) (*httpclient.Response, error) {
			return &httpclient.Response{
				Tables: []httpclient.Table{
					{
						Unixtime:    1612020717,
						Temperature: 25.4,
						Humidity:    65.7,
						Illuminance: 234,
					},
					{
						Unixtime:    1612020733,
						Temperature: 21.9,
						Humidity:    43.0,
						Illuminance: 406,
					},
				},
			}, nil
		},
	})

	m, err := g.GetMetrics(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	want := domain.MetricsList([]domain.Metrics{
		{
			Time:        time.Date(2021, 1, 30, 15, 31, 57, 0, time.UTC),
			Temperature: 25.4,
			Humidity:    65.7,
			Illuminance: 234,
		},
		{
			Time:        time.Date(2021, 1, 30, 15, 32, 13, 0, time.UTC),
			Temperature: 21.9,
			Humidity:    43.0,
			Illuminance: 406,
		},
	})

	if diff := cmp.Diff(want, m); diff != "" {
		t.Errorf("returned metrics differs\n%s", diff)
	}
}

func TestGetMetrics_when_empty_response(t *testing.T) {
	g := mflight.NewMetricsGetter(&stubClient{
		MockGetSensorMonitor: func(context.Context) (*httpclient.Response, error) {
			return &httpclient.Response{}, nil
		},
	})

	m, err := g.GetMetrics(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(domain.MetricsList{}, m); diff != "" {
		t.Errorf("returned metrics is not empty\n%s", diff)
	}
}

func TestGetMetrics_when_request_failure(t *testing.T) {
	c := &stubClient{
		MockGetSensorMonitor: func(context.Context) (*httpclient.Response, error) {
			return &httpclient.Response{}, errors.New("test")
		},
	}

	g := mflight.NewMetricsGetter(c)

	m, err := g.GetMetrics(context.Background())
	if err == nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(domain.MetricsList{}, m); diff != "" {
		t.Errorf("returned metrics is not empty\n%s", diff)
	}
}
