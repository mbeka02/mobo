package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	// Apply global middleware
	r.Use(httprate.LimitByIP(100, time.Minute))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	// logs requests
	r.Use(middleware.Logger)
	// catches panics in the handlers and returns a 500 instead of crashing the server
	r.Use(middleware.Recoverer)
	// extracts the real client IP from the headers even when behind a proxy
	r.Use(middleware.RealIP)
	// add a unique request ID for each request
	r.Use(middleware.RequestID)
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/", s.testHandler)
		// Health Check
		r.Get("/health", s.healthHandler)
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
	stats := s.store.Health()
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(stats)
}
