package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/mbeka02/ticketing-service/config"
	"github.com/mbeka02/ticketing-service/internal/auth"
	"github.com/mbeka02/ticketing-service/internal/server"
	"github.com/mbeka02/ticketing-service/pkg/logger"
	"go.uber.org/zap"
)

func HandleGracefulShutdown(srv *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	logger.Info("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server exiting")
	done <- true
}

func setupServer() (*http.Server, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	if err := logger.Init(cfg.ServerEnv); err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	srv, err := server.NewServer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to setup server: %w", err)
	}

	callbackURL := fmt.Sprintf("http://localhost%s/api/v1/auth/google/callback", cfg.ServerPort)
	auth.NewAuth(callbackURL)

	logger.Info("completed server setup",
		zap.String("port", cfg.ServerPort),
		zap.String("env", cfg.ServerEnv),
	)

	return srv, nil
}

func main() {
	srv, err := setupServer()
	if err != nil {
		logger.Fatal("failed to setup server", zap.Error(err))
	}

	// Ensure logs are flushed on exit
	defer logger.Sync()

	done := make(chan bool, 1)
	go HandleGracefulShutdown(srv, done)

	logger.Info("starting HTTP server")
	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Fatal("http server error", zap.Error(err))
	}

	<-done
	logger.Info("graceful shutdown complete")
}
