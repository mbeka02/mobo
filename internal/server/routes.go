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

	// 1. Recoverer should be first to catch panics from all other middleware
	r.Use(middleware.Recoverer)

	// RealIP extracts real client IP
	r.Use(middleware.RealIP)

	// Request ID middleware
	r.Use(customMiddleware.RequestIDMiddleware)

	// Logging middleware (uses the request ID from step 3)
	r.Use(customMiddleware.LoggingMiddleware)

	// 5. Rate limiting
	r.Use(httprate.LimitByIP(100, time.Minute))

	// 6. CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/", s.testHandler)
		// Health Check
		r.Get("/health", s.healthHandler)
		// Auth
		r.Get("/auth/{provider}", s.handlers.Auth.BeginAuthHandler)
		r.Get("/auth/{provider}/callback", s.handlers.Auth.GetAuthCallbackHandler)
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
