package application_test

import (
	"context"
	"errors"
	"mflight-api/application"
	"mflight-api/domain"
	"mflight-api/domain/mock_domain"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

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
		m    []domain.Metrics
		want domain.Metrics
		err  error
	}{
		{[]domain.Metrics{a, b, c}, c, nil},
		{[]domain.Metrics{a, b}, b, nil},
		{[]domain.Metrics{a}, a, nil},
		{[]domain.Metrics{}, domain.Metrics{}, errors.New("empty metrics")},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockSensor := mock_domain.NewMockSensor(ctrl)
		mockSensor.EXPECT().GetMetrics(ctx).Return(
			domain.TimeSeriesMetrics(tt.m),
			nil,
		)

		collector := application.NewMetricsCollector(mockSensor)

		got, err := collector.CollectLatestMetrics(ctx)

		if err != tt.err && err.Error() != tt.err.Error() {
			t.Errorf("err differs\n got=%v\nwant=%v\n", err, tt.err)
		}
		if diff := cmp.Diff(got, tt.want); diff != "" {
			t.Errorf("returned metrics differs\n%s", diff)
		}
	}
}

func TestCollectTimeSeriesMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	want := domain.TimeSeriesMetrics([]domain.Metrics{a, b, c})

	mockSensor := mock_domain.NewMockSensor(ctrl)
	mockSensor.EXPECT().GetMetrics(ctx).Return(want, nil)

	collector := application.NewMetricsCollector(mockSensor)

	got, err := collector.CollectTimeSeriesMetrics(ctx)

	if err != nil {
		t.Errorf("error: %v\n", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("returned metrics differs\n%s", diff)
	}
}
