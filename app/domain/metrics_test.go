package domain_test

import (
	"mflight-api/app/domain"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestMetricsList(t *testing.T) {
	t.Run("Last()", func(t *testing.T) {
		t.Run("should return last element", func(t *testing.T) {
			l := domain.MetricsList([]domain.Metrics{
				{
					Time:        time.Time{},
					Temperature: 25,
					Humidity:    62,
					Illuminance: 412,
				},
				{
					Time:        time.Time{},
					Temperature: 24,
					Humidity:    60,
					Illuminance: 400,
				},
			})

			want := domain.Metrics{
				Time:        time.Time{},
				Temperature: 24,
				Humidity:    60,
				Illuminance: 400,
			}

			got, err := l.Last()

			if err != nil {
				t.Fatalf("error must not occur. %v", err)
			}

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("differ. \n%s\n", diff)
			}
		})

		t.Run("should return error when list has no element", func(t *testing.T) {
			l := domain.MetricsList([]domain.Metrics{})

			_, err := l.Last()

			if err == nil {
				t.Fatal("error must occur.")
			}
		})
	})
}
