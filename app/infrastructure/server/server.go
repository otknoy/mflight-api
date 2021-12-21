package server

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
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

		zap.L().Info("shutdown gracefully...")
		if err := s.Shutdown(shutdownCtx); err != nil {
			zap.L().Error("shutdown error", zap.Error(err))
		}

		close(idleConnsClosed)
	}()

	err := s.ListenAndServe()

	<-idleConnsClosed

	return err
}
