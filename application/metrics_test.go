package application_test

import (
	"context"
	"mflight-api/application"
	"mflight-api/domain"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type stubSensor struct{}

func (s *stubSensor) GetMetrics(ctx context.Context) (domain.TimeSeriesMetrics, error) {
	return domain.TimeSeriesMetrics([]domain.Metrics{
		{
			Temperature: domain.Temperature(18.0),
			Humidity:    domain.Humidity(45.0),
			Illuminance: domain.Illuminance(300),
		},
	}), nil
}

func TestCollectMetrics(t *testing.T) {
	c := application.NewMetricsCollector(&stubSensor{})

	m, _ := c.CollectMetrics(context.Background())

	want := domain.TimeSeriesMetrics([]domain.Metrics{
		{
			Temperature: 18.0,
			Humidity:    45.0,
			Illuminance: 300,
		},
	})

	if diff := cmp.Diff(want, m); diff != "" {
		t.Errorf("returned metrics differs\n%s", diff)
	}
}
