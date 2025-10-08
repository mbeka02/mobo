package main

import (
	"log"

	"github.com/mbeka02/ticketing-service/config"
	"github.com/mbeka02/ticketing-service/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	log.Printf("Starting server on port %s (env: %s)", cfg.ServerPort, cfg.ServerEnv)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

