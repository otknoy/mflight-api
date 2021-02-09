package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type HTTPServer struct {
	srv *http.Server
}

func (s *HTTPServer) ListenAndServe(idleConnsClosed chan<- struct{}) {
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := s.srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}

		close(idleConnsClosed)
	}()
	if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}

func NewServer(handlers map[string]http.Handler, port int) *HTTPServer {
	mux := http.NewServeMux()
	for path, handler := range handlers {
		mux.Handle(path, handler)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return &HTTPServer{srv}
}
