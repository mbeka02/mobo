package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	// Server config
	ServerPort         string        `mapstructure:"SERVER_PORT"`
	ServerEnv          string        `mapstructure:"SERVER_ENV"`
	BaseURL            string        `mapstructure:"BASE_URL"`
	FrontendURL        string        `mapstructure:"FRONTEND_URL"`
	ServerReadTimeout  time.Duration `mapstructure:"SERVER_READTIMEOUT"`
	ServerWriteTimeout time.Duration `mapstructure:"SERVER_WRITETIMEOUT"`
	ServerIdleTimeout  time.Duration `mapstructure:"SERVER_IDLETIMEOUT"`

	// Auth config
	SymmetricKey         string        `mapstructure:"SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`

	// Database config
	DatabaseURI             string        `mapstructure:"DATABASE_URI"`
	DatabaseMaxConnections  int           `mapstructure:"DATABASE_MAXCONNECTIONS"`
	DatabaseMinConnections  int           `mapstructure:"DATABASE_MINCONNECTIONS"`
	DatabaseMaxConnLifetime time.Duration `mapstructure:"DATABASE_MAXCONNLIFETIME"`

	// Async email worker configuration. These are only required by cmd/email-worker.
	RabbitMQURL         string        `mapstructure:"RABBITMQ_URL"`
	EmailProvider       string        `mapstructure:"EMAIL_PROVIDER"`
	EmailFrom           string        `mapstructure:"EMAIL_FROM"`
	ResendAPIKey        string        `mapstructure:"RESEND_API_KEY"`
	EmailWorkerCount    int           `mapstructure:"EMAIL_WORKER_COUNT"`
	EmailWorkerPrefetch int           `mapstructure:"EMAIL_WORKER_PREFETCH"`
	EmailSendTimeout    time.Duration `mapstructure:"EMAIL_SEND_TIMEOUT"`
	OutboxPollInterval  time.Duration `mapstructure:"OUTBOX_POLL_INTERVAL"`
	OutboxLeaseDuration time.Duration `mapstructure:"OUTBOX_LEASE_DURATION"`
	OutboxBatchSize     int           `mapstructure:"OUTBOX_BATCH_SIZE"`
}

type DatabaseConfig struct {
	URI             string
	MaxConnections  int
	MinConnections  int
	MaxConnLifetime time.Duration
}

func (c *Config) GetDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		URI:             c.DatabaseURI,
		MaxConnections:  c.DatabaseMaxConnections,
		MinConnections:  c.DatabaseMinConnections,
		MaxConnLifetime: c.DatabaseMaxConnLifetime,
	}
}

func Load() (*Config, error) {
	v := viper.New()
	setDefaults(v)

	// Try to load .env file (optional in production)
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	envFileErr := v.ReadInConfig()

	// Auto-detect environment if SERVER_ENV is not explicitly set
	if os.Getenv("SERVER_ENV") == "" && !v.IsSet("SERVER_ENV") {
		if envFileErr == nil {
			// .env file exists — we're in development
			v.Set("SERVER_ENV", "development")
		} else {
			// No .env file — assume production
			v.Set("SERVER_ENV", "production")
		}
	}

	// Explicitly bind environment variables
	bindEnvVars(v)

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// Ensure port format
	if cfg.ServerPort != "" && cfg.ServerPort[0] != ':' {
		cfg.ServerPort = ":" + cfg.ServerPort
	}

	return &cfg, nil
}

func bindEnvVars(v *viper.Viper) {
	// Server vars
	envVars := []string{
		"SERVER_PORT",
		"SERVER_ENV",
		"SERVER_READTIMEOUT",
		"SERVER_WRITETIMEOUT",
		"SERVER_IDLETIMEOUT",
		"DATABASE_URI",
		"DATABASE_MAXCONNECTIONS",
		"DATABASE_MINCONNECTIONS",
		"DATABASE_MAXCONNLIFETIME",
		"SYMMETRIC_KEY",
		"ACCESS_TOKEN_DURATION",
		"REFRESH_TOKEN_DURATION",
		"BASE_URL",
		"FRONTEND_URL",
		"RABBITMQ_URL",
		"EMAIL_PROVIDER",
		"EMAIL_FROM",
		"RESEND_API_KEY",
		"EMAIL_WORKER_COUNT",
		"EMAIL_WORKER_PREFETCH",
		"EMAIL_SEND_TIMEOUT",
		"OUTBOX_POLL_INTERVAL",
		"OUTBOX_LEASE_DURATION",
		"OUTBOX_BATCH_SIZE",
	}

	for _, envVar := range envVars {
		if val := os.Getenv(envVar); val != "" {
			v.Set(envVar, val)
		}
	}
}

func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("SERVER_PORT", "3000")
	v.SetDefault("BASE_URL", "http://localhost:3000")
	v.SetDefault("FRONTEND_URL", "http://localhost:5173")
	v.SetDefault("SERVER_READTIMEOUT", 45*time.Second)
	v.SetDefault("SERVER_WRITETIMEOUT", 30*time.Second)
	v.SetDefault("SERVER_IDLETIMEOUT", time.Minute)

	// Database defaults
	v.SetDefault("DATABASE_MAXCONNECTIONS", 25)
	v.SetDefault("DATABASE_MINCONNECTIONS", 5)
	v.SetDefault("DATABASE_MAXCONNLIFETIME", 30*time.Minute)

	// Auth defaults
	v.SetDefault("ACCESS_TOKEN_DURATION", 15*time.Minute)
	v.SetDefault("REFRESH_TOKEN_DURATION", 7*24*time.Hour)

	// Async email defaults. Credentials and broker URL intentionally have no defaults.
	v.SetDefault("EMAIL_PROVIDER", "resend")
	v.SetDefault("EMAIL_WORKER_COUNT", 2)
	v.SetDefault("EMAIL_WORKER_PREFETCH", 10)
	v.SetDefault("EMAIL_SEND_TIMEOUT", 10*time.Second)
	v.SetDefault("OUTBOX_POLL_INTERVAL", 5*time.Second)
	v.SetDefault("OUTBOX_LEASE_DURATION", time.Minute)
	v.SetDefault("OUTBOX_BATCH_SIZE", 25)
}

// ValidateEmailWorker checks configuration that is intentionally optional for the HTTP API.
func (c *Config) ValidateEmailWorker() error {
	if c.RabbitMQURL == "" {
		return fmt.Errorf("RABBITMQ_URL is required")
	}
	if c.EmailFrom == "" {
		return fmt.Errorf("EMAIL_FROM is required")
	}
	if c.EmailProvider != "resend" {
		return fmt.Errorf("EMAIL_PROVIDER must currently be resend")
	}
	if c.ResendAPIKey == "" {
		return fmt.Errorf("RESEND_API_KEY is required for resend")
	}
	if c.EmailWorkerCount < 1 {
		return fmt.Errorf("EMAIL_WORKER_COUNT must be at least 1")
	}
	if c.EmailWorkerPrefetch < 1 {
		return fmt.Errorf("EMAIL_WORKER_PREFETCH must be at least 1")
	}
	if c.OutboxPollInterval <= 0 || c.OutboxLeaseDuration <= 0 {
		return fmt.Errorf("outbox durations must be positive")
	}
	if c.OutboxBatchSize < 1 {
		return fmt.Errorf("OUTBOX_BATCH_SIZE must be at least 1")
	}
	return nil
}

func (c *Config) Validate() error {
	// Required fields
	if c.DatabaseURI == "" {
		return fmt.Errorf("DATABASE_URI is required but not set")
	}

	if len(c.SymmetricKey) < 32 {
		return fmt.Errorf("SYMMETRIC_KEY must be at least 32 characters")
	}

	// Constraints
	if c.DatabaseMaxConnections < c.DatabaseMinConnections {
		return fmt.Errorf("DATABASE_MAXCONNECTIONS must be >= DATABASE_MINCONNECTIONS")
	}

	if c.DatabaseMaxConnections < 1 {
		return fmt.Errorf("DATABASE_MAXCONNECTIONS must be at least 1")
	}

	// Validate environment
	validEnvs := map[string]bool{"development": true, "staging": true, "production": true}
	if !validEnvs[c.ServerEnv] {
		return fmt.Errorf("SERVER_ENV must be one of: development, staging, production")
	}

	return nil
}

// IsProduction returns true if running in production
func (c *Config) IsProduction() bool {
	return c.ServerEnv == "production"
}
