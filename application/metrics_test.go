package application_test

import (
	"context"
	"errors"
	"mflight-api/application"
	"mflight-api/domain"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type mockSensor struct {
	domain.MetricsGetter
	MockGetMetrics func(ctx context.Context) ([]domain.Metrics, error)
}

func (s *mockSensor) GetMetrics(ctx context.Context) ([]domain.Metrics, error) {
	return s.MockGetMetrics(ctx)
}

var (
	a = domain.Metrics{
		Time:        time.Date(2021, 1, 31, 13, 8, 5, 0, time.UTC),
		Temperature: domain.Temperature(18.0),
		Humidity:    domain.Humidity(45.0),
		Illuminance: domain.Illuminance(300),
	}
	b = domain.Metrics{
		Time:        time.Date(2021, 1, 31, 13, 9, 5, 0, time.UTC),
		Temperature: domain.Temperature(19.0),
		Humidity:    domain.Humidity(46.0),
		Illuminance: domain.Illuminance(301),
	}
	c = domain.Metrics{
		Time:        time.Date(2021, 1, 31, 15, 9, 5, 0, time.UTC),
		Temperature: domain.Temperature(19.1),
		Humidity:    domain.Humidity(47.0),
		Illuminance: domain.Illuminance(500),
	}
)

func TestCollectLatestMetrics(t *testing.T) {
	tests := []struct {
		m       []domain.Metrics
		err     error
		want    domain.Metrics
		wantErr error
	}{
		{
			[]domain.Metrics{a, b, c}, nil,
			c, nil,
		},
		{
			[]domain.Metrics{a, b}, nil,
			b, nil,
		},
		{
			[]domain.Metrics{a}, nil,
			a, nil,
		},
		{
			[]domain.Metrics{}, errors.New("error"),
			domain.Metrics{}, errors.New("error"),
		},
		{
			[]domain.Metrics{}, nil,
			domain.Metrics{}, errors.New("no metrics"),
		},
	}

	for _, tt := range tests {
		testCtx := context.Background()

		collector := application.NewMetricsCollector(&mockSensor{
			MockGetMetrics: func(ctx context.Context) ([]domain.Metrics, error) {
				if ctx != testCtx {
					t.Fail()
				}
				return tt.m, tt.err
			},
		})

		got, err := collector.CollectLatestMetrics(testCtx)

		if err != tt.wantErr && err.Error() != tt.wantErr.Error() {
			t.Errorf("err differs\n got=%v\nwant=%v\n", err, tt.err)
		}
		if diff := cmp.Diff(got, tt.want); diff != "" {
			t.Errorf("returned metrics differs\n%s", diff)
		}
	}
}

func TestCollectTimeSeriesMetrics(t *testing.T) {
	ctx := context.Background()

	want := []domain.Metrics{a, b, c}

	collector := application.NewMetricsCollector(&mockSensor{
		MockGetMetrics: func(ctx context.Context) ([]domain.Metrics, error) {
			return want, nil
		},
	})

	got, err := collector.CollectMetricsList(ctx)

	if err != nil {
		t.Errorf("error: %v\n", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("returned metrics differs\n%s", diff)
	}
}
