package middleware_test

import (
	"errors"
	"mflight-api/infrastructure/prometheus/middleware"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type stubRoundTripper struct {
	StubRoundTrip func(*http.Request) (*http.Response, error)
}

func (rt *stubRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return rt.StubRoundTrip(r)
}

func TestRoundTrip(t *testing.T) {
	in := &http.Request{Method: "GET", URL: &url.URL{Host: "example.com:56002", Path: "/foo"}}
	want := &http.Response{Status: "200 OK"}

	h := middleware.NewRoundTripperMetricsMiddleware(
		&stubRoundTripper{
			func(r *http.Request) (*http.Response, error) {
				if r == in {
					return want, nil
				}

				return &http.Response{}, errors.New("error")
			},
		},
	)

	got, err := h.RoundTrip(in)

	if err != nil {
		t.Errorf("err should not be nil.\n%v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("response differs.\n%v", diff)
	}
}
