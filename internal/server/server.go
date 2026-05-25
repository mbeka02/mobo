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
	Auth      *handler.AuthHandler
	Movie     *handler.MovieHandler
	Showtime  *handler.ShowtimeHandler
	Venue     *handler.VenueHandler
	Dashboard *handler.DashboardHandler
}

type Services struct {
	Auth      service.AuthService
	Movie     service.MovieService
	Showtime  service.ShowtimeService
	Venue     service.VenueService
	Dashboard service.DashboardService
}

type Repositories struct {
	Auth      repository.AuthRepository
	Movie     repository.MovieRepository
	Showtime  repository.ShowtimeRepository
	Venue     repository.VenueRepository
	Dashboard repository.DashboardRepository
}

type Server struct {
	store      *database.Store
	config     *config.Config
	handlers   *Handlers
	tokenMaker auth.Maker
}

func initRepositories(store *database.Store) *Repositories {
	return &Repositories{
		Auth:      repository.NewAuthRepository(store),
		Movie:     repository.NewMovieRepository(store),
		Showtime:  repository.NewShowtimeRepository(store),
		Venue:     repository.NewVenueRepository(store),
		Dashboard: repository.NewDashboardRepository(store),
	}
}

func initServices(repos *Repositories) *Services {
	return &Services{
		Auth:      service.NewAuthService(repos.Auth),
		Movie:     service.NewMovieService(repos.Movie),
		Showtime:  service.NewShowtimeService(repos.Showtime),
		Venue:     service.NewVenueService(repos.Venue),
		Dashboard: service.NewDashboardService(repos.Dashboard),
	}
}

func initHandlers(services *Services, maker auth.Maker, cfg *config.Config) *Handlers {
	return &Handlers{
		Auth:      handler.NewAuthHandler(services.Auth, maker, cfg.IsProduction(), cfg.AccessTokenDuration, cfg.RefreshTokenDuration, cfg.FrontendURL),
		Movie:     handler.NewMovieHandler(services.Movie),
		Showtime:  handler.NewShowtimeHandler(services.Showtime),
		Venue:     handler.NewVenueHandler(services.Venue),
		Dashboard: handler.NewDashboardHandler(services.Dashboard),
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

