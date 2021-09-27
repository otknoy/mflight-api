package server

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type GracefulShutdownServer struct {
	http.Server
}

func (s *GracefulShutdownServer) ListenAndServeGracefully() error {
	idleConnsClosed := make(chan struct{})

	go func() {
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Println("shutdown gracefully...")
		if err := s.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown error: %v", err)
		}

		close(idleConnsClosed)
	}()

	err := s.ListenAndServe()

	<-idleConnsClosed

	return err
}
