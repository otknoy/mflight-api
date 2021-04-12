package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

// GracefulShutdownServer is a wrapper of http.Server to gracefully shutdown
type GracefulShutdownServer struct {
	http.Server
}

func (s *GracefulShutdownServer) ListenAndServeWithGracefulShutdown(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("shutdown gracefully...")
	err := s.Shutdown(shutdownCtx)
	if err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}

	return err
}

func NewServer(mux *http.ServeMux, port int) *GracefulShutdownServer {
	return &GracefulShutdownServer{
		http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}
}
