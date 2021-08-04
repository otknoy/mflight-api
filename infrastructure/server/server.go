package server

import (
	"context"
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

func (s *GracefulShutdownServer) ListenAndServeWithGracefulShutdown() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	done := make(chan struct{})

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Println("shutdown gracefully...")
		if err := s.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown error: %v", err)
		}

		close(done)
	}()

	err := s.ListenAndServe()

	<-done

	return err
}
