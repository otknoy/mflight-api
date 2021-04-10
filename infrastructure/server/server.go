package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// GracefulShutdownServer is a wrapper of http.Server to gracefully shutdown
type GracefulShutdownServer struct {
	http.Server
}

func (s *GracefulShutdownServer) ListenAndServeWithGracefulShutdown() <-chan struct{} {
	idleConnsClosed := make(chan struct{})

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

		log.Printf("signal: %v", <-sig)

		if err := s.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}

		close(idleConnsClosed)
	}()
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	return idleConnsClosed
}

func NewServer(mux *http.ServeMux, port int) *GracefulShutdownServer {
	return &GracefulShutdownServer{
		http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}
}