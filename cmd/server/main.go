package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/mbeka02/ticketing-service/config"
	"github.com/mbeka02/ticketing-service/internal/auth"
	"github.com/mbeka02/ticketing-service/internal/server"
)

func HandleGracefulShutdown(srv *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	log.Println("shutting down gracefully, press Ctrl+C again to force")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdown with error: %v", err)
	}
	log.Println("server exiting")
	done <- true
}

func setupServer() (*http.Server, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration:(%w)", err)
	}
	srv, err := server.NewServer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to setup server(%w)", err)
	}
	auth.NewAuth()
	log.Printf("Completed server setup , starting server on port %s (env: %s)", cfg.ServerPort, cfg.ServerEnv)
	return srv, nil
}

func main() {
	srv, err := setupServer()
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan bool, 1)
	go HandleGracefulShutdown(srv, done)

	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}
	<-done
	log.Println("Graceful shutdown complete")
}
