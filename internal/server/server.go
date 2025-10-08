package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mbeka02/ticketing-service/config"
	"github.com/mbeka02/ticketing-service/internal/database"
)

type Server struct {
	store  *database.Store
	config *config.Config
}

func NewServer(cfg *config.Config) (*http.Server, error) {
	// Create database store with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	store, err := database.NewStore(ctx, cfg.GetDatabaseConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Test database connection
	if err := store.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	srv := &Server{
		store:  store,
		config: cfg,
	}

	return &http.Server{
		Handler:      srv.RegisterRoutes(),
		Addr:         cfg.ServerPort,
		IdleTimeout:  cfg.ServerIdleTimeout,
		ReadTimeout:  cfg.ServerReadTimeout,
		WriteTimeout: cfg.ServerWriteTimeout,
	}, nil
}

