package fixedtime_test

import (
	"mflight-api/app/infrastructure/fixedtime"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestNow(t *testing.T) {
	a := time.Now()

	got := fixedtime.Now()

	b := time.Now()

	if a.After(got) && got.After(b) {
		t.Error("got should be in between a and b.")
	}
}

func TestSet(t *testing.T) {
	tests := []time.Time{
		time.Date(1988, 12, 16, 23, 21, 55, 0, time.UTC),
		time.Date(2021, 1, 31, 23, 21, 56, 0, time.UTC),
		time.Date(2038, 12, 16, 23, 21, 55, 0, time.UTC),
	}

	for _, tt := range tests {
		fixedtime.Set(func() time.Time {
			return tt
		})

		got := fixedtime.Now()

		if diff := cmp.Diff(tt, got); diff != "" {
			t.Errorf("got should be return value of set function\n%v", diff)
		}
	}
}
