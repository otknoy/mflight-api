package middleware

import "net/http"

type statusRecoder struct {
	status int
	w      http.ResponseWriter
}

func newStatusRecoder(w http.ResponseWriter) *statusRecoder {
	return &statusRecoder{
		status: 0,
		w:      w,
	}
}

func (r *statusRecoder) Header() http.Header {
	return r.w.Header()
}

func (r *statusRecoder) Write(bytes []byte) (int, error) {
	return r.w.Write(bytes)
}

func (r *statusRecoder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.w.WriteHeader(statusCode)
}
