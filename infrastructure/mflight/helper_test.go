package mflight_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func NewStubServer(t *testing.T) *httptest.Server {
	t.Helper()

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if path := r.URL.Path; path != "/SensorMonitorV2.xml" {
			t.Fatal(path)
		}

		if qs := r.URL.RawQuery; qs != "x-KEY_MOBILE_ID=test-mobile-id&x-KEY_UPDATE_DATE=" {
			t.Fatal(qs)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(
			w,
			`
      <db>
        <table id="67243">
          <time>202101030000</time>
          <unixtime>1609599600</unixtime>
          <temp>22.0</temp>
          <humi>43.3</humi>
          <illu>405</illu>
        </table>
        <table id="67244">
          <time>202101030005</time>
          <unixtime>1609599900</unixtime>
          <temp>21.9</temp>
          <humi>43.0</humi>
          <illu>406</illu>
        </table>
      </db>`,
		)
	}))

	return s
}
