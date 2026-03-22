package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mbeka02/ticketing-service/config"
	"github.com/mbeka02/ticketing-service/internal/auth"
	"github.com/mbeka02/ticketing-service/internal/database"
	"github.com/mbeka02/ticketing-service/internal/server/handler"
	"github.com/mbeka02/ticketing-service/internal/server/repository"
	"github.com/mbeka02/ticketing-service/internal/server/service"
	"github.com/mbeka02/ticketing-service/pkg/logger"
	"go.uber.org/zap"
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
	store      *database.Store
	config     *config.Config
	handlers   *Handlers
	tokenMaker auth.Maker
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

func initHandlers(services *Services, maker auth.Maker, cfg *config.Config) *Handlers {
	return &Handlers{
		Auth: handler.NewAuthHandler(services.Auth, maker, cfg.IsProduction(), cfg.AccessTokenDuration, cfg.RefreshTokenDuration),
	}
}

func NewServer(cfg *config.Config) (*http.Server, error) {
	logger.Info("initializing server")

	// Create database store with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("connecting to database")

	store, err := database.NewStore(ctx, cfg.GetDatabaseConfig())
	if err != nil {
		logger.Error("failed to initialize database", zap.Error(err))
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Test database connection
	logger.Debug("pinging database")
	if err := store.Ping(ctx); err != nil {
		logger.Error("failed to ping database", zap.Error(err))
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("database connection established successfully")

	// Create JWT token maker
	tokenMaker, err := auth.NewJWTMaker(cfg.SymmetricKey)
	if err != nil {
		logger.Error("failed to create token maker", zap.Error(err))
		return nil, fmt.Errorf("failed to create token maker: %w", err)
	}

	logger.Info("token maker initialized successfully")

	repos := initRepositories(store)
	services := initServices(repos)
	handlers := initHandlers(services, tokenMaker, cfg)

	srv := &Server{
		store:      store,
		config:     cfg,
		handlers:   handlers,
		tokenMaker: tokenMaker,
	}

	return &http.Server{
		Handler:      srv.RegisterRoutes(),
		Addr:         cfg.ServerPort,
		IdleTimeout:  cfg.ServerIdleTimeout,
		ReadTimeout:  cfg.ServerReadTimeout,
		WriteTimeout: cfg.ServerWriteTimeout,
	}, nil
}

