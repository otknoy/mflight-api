package middleware_test

import (
	"fmt"
	"mflight-api/app/infrastructure/prometheus/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type stubHandler struct {
	StubServeHTTP func(w http.ResponseWriter, r *http.Request)
}

func (h *stubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.StubServeHTTP(w, r)
}

func TestServeHTTP(t *testing.T) {
	h := middleware.InstrumentHandlerMetrics(
		&stubHandler{
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, "stub-response")
			},
		},
	)

	req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
	got := httptest.NewRecorder()

	h.ServeHTTP(got, req)

	if diff := cmp.Diff(got.Code, http.StatusOK); diff != "" {
		t.Errorf("invalid response status: %v\n", diff)
	}
	if diff := cmp.Diff(got.Body.String(), "stub-response"); diff != "" {
		t.Errorf("invalid response body: %v\n", diff)
	}
}
