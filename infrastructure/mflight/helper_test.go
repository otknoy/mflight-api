package mflight_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func NewStubServer(t *testing.T, response string) *httptest.Server {
	t.Helper()

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if path := r.URL.Path; path != "/SensorMonitorV2.xml" {
			t.Fatal(path)
		}

		if qs := r.URL.RawQuery; qs != "x-KEY_MOBILE_ID=test-mobile-id&x-KEY_UPDATE_DATE=" {
			t.Fatal(qs)
		}

		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, response)
	}))

	return s
}
