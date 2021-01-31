package application_test

import (
	"context"
	"mflight-api/application"
	"mflight-api/domain"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type stubSensor struct{}

func (s *stubSensor) GetMetrics(ctx context.Context) (domain.TimeSeriesMetrics, error) {
	return domain.TimeSeriesMetrics([]domain.Metrics{
		{
			Time:        time.Date(2021, 1, 31, 13, 8, 5, 0, time.UTC),
			Temperature: domain.Temperature(18.0),
			Humidity:    domain.Humidity(45.0),
			Illuminance: domain.Illuminance(300),
		},
		{
			Time:        time.Date(2021, 1, 31, 13, 9, 5, 0, time.UTC),
			Temperature: domain.Temperature(19.0),
			Humidity:    domain.Humidity(46.0),
			Illuminance: domain.Illuminance(301),
		},
	}), nil
}

func TestCollectLatestMetrics(t *testing.T) {
	c := application.NewMetricsCollector(&stubSensor{})

	got, _ := c.CollectLatestMetrics(context.Background())

	want := domain.Metrics{
		Time:        time.Date(2021, 1, 31, 13, 9, 5, 0, time.UTC),
		Temperature: 19.0,
		Humidity:    46.0,
		Illuminance: 301,
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("returned metrics differs\n%s", diff)
	}
}

func TestCollectTimeSeriesMetrics(t *testing.T) {
	c := application.NewMetricsCollector(&stubSensor{})

	got, _ := c.CollectTimeSeriesMetrics(context.Background())

	want := domain.TimeSeriesMetrics([]domain.Metrics{
		{
			Time:        time.Date(2021, 1, 31, 13, 8, 5, 0, time.UTC),
			Temperature: 18.0,
			Humidity:    45.0,
			Illuminance: 300,
		},
		{
			Time:        time.Date(2021, 1, 31, 13, 9, 5, 0, time.UTC),
			Temperature: 19.0,
			Humidity:    46.0,
			Illuminance: 301,
		},
	})

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("returned metrics differs\n%s", diff)
	}
}
