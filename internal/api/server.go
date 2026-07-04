package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mbeka02/ticketing-service/config"
	"github.com/mbeka02/ticketing-service/internal/analytics"
	"github.com/mbeka02/ticketing-service/internal/auth"
	"github.com/mbeka02/ticketing-service/internal/movie"
	"github.com/mbeka02/ticketing-service/internal/postgres"
	"github.com/mbeka02/ticketing-service/internal/showtime"
	"github.com/mbeka02/ticketing-service/internal/user"
	"github.com/mbeka02/ticketing-service/internal/venue"
	"github.com/mbeka02/ticketing-service/pkg/logger"
	"go.uber.org/zap"
)

// Handlers groups all HTTP handlers.
type Handlers struct {
	User      *UserHandler
	Movie     *MovieHandler
	Showtime  *ShowtimeHandler
	Venue     *VenueHandler
	Analytics *AnalyticsHandler
}

// Server holds dependencies for the HTTP server.
type Server struct {
	store      *postgres.Store
	config     *config.Config
	handlers   *Handlers
	tokenMaker auth.Maker
}

// NewServer creates and configures a new HTTP server.
func NewServer(cfg *config.Config) (*http.Server, error) {
	logger.Info("initializing server")

	// Create database store with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("connecting to database")

	store, err := postgres.NewStore(ctx, cfg.GetDatabaseConfig())
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

	// Initialize repositories (postgres adapters)
	userRepo := postgres.NewUserRepository(store)
	movieRepo := postgres.NewMovieRepository(store)
	showtimeRepo := postgres.NewShowtimeRepository(store)
	venueRepo := postgres.NewVenueRepository(store)
	analyticsRepo := postgres.NewAnalyticsRepository(store)

	// Initialize domain services
	userSvc := user.NewService(userRepo)
	movieSvc := movie.NewService(movieRepo)
	showtimeSvc := showtime.NewService(showtimeRepo)
	venueSvc := venue.NewService(venueRepo)
	analyticsSvc := analytics.NewService(analyticsRepo)

	// Initialize handlers
	handlers := &Handlers{
		User:      NewUserHandler(userSvc, tokenMaker, cfg.IsProduction(), cfg.AccessTokenDuration, cfg.RefreshTokenDuration, cfg.FrontendURL),
		Movie:     NewMovieHandler(movieSvc),
		Showtime:  NewShowtimeHandler(showtimeSvc),
		Venue:     NewVenueHandler(venueSvc),
		Analytics: NewAnalyticsHandler(analyticsSvc),
	}

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
