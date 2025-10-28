package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mbeka02/ticketing-service/config"
	"github.com/mbeka02/ticketing-service/internal/database"
	"github.com/mbeka02/ticketing-service/internal/server/handler"
	"github.com/mbeka02/ticketing-service/internal/server/repository"
	"github.com/mbeka02/ticketing-service/internal/server/service"
)

type Handlers struct {
	Auth *handler.AuthHandler
}
type Services struct {
	Auth service.AuthService
}

type Repositories struct {
	Auth repository.AuthRepository
}
type Server struct {
	store    *database.Store
	config   *config.Config
	handlers *Handlers
}

func initRepositories(store *database.Store) *Repositories {
	return &Repositories{
		Auth: repository.NewAuthRepository(store),
	}
}

func initServices(repos *Repositories) *Services {
	return &Services{
		Auth: service.NewAuthService(repos.Auth),
	}
}

func initHandlers(services *Services) *Handlers {
	return &Handlers{
		Auth: handler.NewAuthHandler(services.Auth),
	}
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
	repos := initRepositories(store)
	services := initServices(repos)
	handlers := initHandlers(services)
	srv := &Server{
		store:    store,
		config:   cfg,
		handlers: handlers,
	}

	return &http.Server{
		Handler:      srv.RegisterRoutes(),
		Addr:         cfg.ServerPort,
		IdleTimeout:  cfg.ServerIdleTimeout,
		ReadTimeout:  cfg.ServerReadTimeout,
		WriteTimeout: cfg.ServerWriteTimeout,
	}, nil
}
