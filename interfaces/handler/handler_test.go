package handler_test

import (
	"context"
	"errors"
	"mflight-api/domain"
	"mflight-api/interfaces/handler"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type mockMetricsGetter struct {
	MockGetMetrics func(context.Context) (domain.MetricsList, error)
}

func (c *mockMetricsGetter) GetMetrics(ctx context.Context) (domain.MetricsList, error) {
	return c.MockGetMetrics(ctx)
}

func TestServeHTTP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/getSensorMetrics", nil)
	got := httptest.NewRecorder()

	h := handler.NewSensorMetricsHandler(&mockMetricsGetter{
		MockGetMetrics: func(ctx context.Context) (domain.MetricsList, error) {
			return []domain.Metrics{
					{
						Time:        time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						Temperature: domain.Temperature(21.3),
						Humidity:    domain.Humidity(52.4),
						Illuminance: domain.Illuminance(400),
					},
					{
						Time:        time.Date(2021, 1, 1, 0, 0, 59, 0, time.UTC),
						Temperature: domain.Temperature(22.5),
						Humidity:    domain.Humidity(50.2),
						Illuminance: domain.Illuminance(401),
					},
				},
				nil
		},
	})

	h.ServeHTTP(got, req)

	if v := got.Code; v != http.StatusOK {
		t.Errorf("http status: 200, but %v\n", v)
	}

	want := `[{"unixtime":1609459200,"temperature":21.3,"humidity":52.4,"illuminance":400},{"unixtime":1609459259,"temperature":22.5,"humidity":50.2,"illuminance":401}]`
	if v := got.Body.String(); v != want {
		t.Errorf("invalid response json:\nwant=%v\n got=%v\n", want, v)
	}
}

func TestServeHTTP_sensor_error(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/getSensorMetrics", nil)
	got := httptest.NewRecorder()

	h := handler.NewSensorMetricsHandler(&mockMetricsGetter{
		MockGetMetrics: func(ctx context.Context) (domain.MetricsList, error) {
			return domain.MetricsList{}, errors.New("failed to get metrics")
		},
	})

	h.ServeHTTP(got, req)

	if v := got.Code; v != http.StatusInternalServerError {
		t.Errorf("http status: 500, but %v\n", v)
	}
	if diff := cmp.Diff(got.Body.String(), "failed to get metrics\n"); diff != "" {
		t.Errorf("response body differs.\n%v", diff)
	}
}
