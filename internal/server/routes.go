package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	customMiddleware "github.com/mbeka02/ticketing-service/internal/server/middleware"
	"github.com/mbeka02/ticketing-service/pkg/logger"
	"go.uber.org/zap"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	// RealIP extracts real client IP
	r.Use(middleware.RealIP)

	// Request ID middleware
	r.Use(customMiddleware.RequestIDMiddleware)

	// Logging middleware (uses the request ID from step 3)
	r.Use(customMiddleware.LoggingMiddleware)

	// Rate limiting
	r.Use(httprate.LimitByIP(100, time.Minute))

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/", s.testHandler)
		// Health Check
		r.Get("/health", s.healthHandler)

		// OAuth routes
		r.Get("/auth/{provider}", s.handlers.Auth.BeginAuthHandler)
		r.Get("/auth/{provider}/callback", s.handlers.Auth.GetAuthCallbackHandler)

		// Traditional auth routes (public)
		r.Post("/auth/signup", s.handlers.Auth.SignupHandler)
		r.Post("/auth/login", s.handlers.Auth.LoginHandler)
		r.Post("/auth/logout", s.handlers.Auth.LogoutHandler)

		// Public listings
		r.Get("/movies", s.handlers.Movie.ListMoviesPublicHandler)
		r.Get("/movies/{movieId}", s.handlers.Movie.GetMovieHandler)
		r.Get("/movies/{movieId}/showtimes", s.handlers.Showtime.ListShowtimesByMovieHandler)
		r.Get("/venues", s.handlers.Venue.ListVenuesHandler)
		r.Get("/venues/{venueId}", s.handlers.Venue.GetVenueHandler)
		r.Get("/showtimes/{showtimeId}", s.handlers.Showtime.GetShowtimeHandler)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(customMiddleware.AuthMiddleware(s.tokenMaker, s.config.IsProduction(), s.config.AccessTokenDuration, s.config.RefreshTokenDuration))

			r.Get("/me", s.handlers.Auth.GetCurrentUser)

			// Admin only routes
			r.Group(func(r chi.Router) {
				r.Use(customMiddleware.AdminMiddleware)

				// Admin Movies
				r.Get("/admin/movies", s.handlers.Movie.ListMoviesAdminHandler)
				r.Post("/admin/movies", s.handlers.Movie.AddMovieHandler)
				r.Patch("/admin/movies/{movieId}", s.handlers.Movie.UpdateMovieHandler)
				r.Delete("/admin/movies/{movieId}", s.handlers.Movie.DeleteMovieHandler)

				// Admin Showtimes
				r.Get("/admin/showtimes", s.handlers.Showtime.ListShowtimesAdminHandler)
				r.Post("/admin/showtimes", s.handlers.Showtime.CreateShowtimeHandler)
				r.Patch("/admin/showtimes/{showtimeId}", s.handlers.Showtime.UpdateShowtimeHandler)
				r.Delete("/admin/showtimes/{showtimeId}", s.handlers.Showtime.DeleteShowtimeHandler)

				// Admin Venues
				r.Post("/admin/venues", s.handlers.Venue.CreateVenueHandler)
			})
		})
	})

	return r
}

func (s *Server) testHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Mobo API"
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger.DebugCtx(ctx, "health check requested")

	stats := s.store.Health()

	logger.InfoCtx(ctx, "health check completed",
		zap.Any("stats", stats),
	)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(stats)
}
