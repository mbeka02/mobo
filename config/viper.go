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
	ServerReadTimeout  time.Duration `mapstructure:"SERVER_READTIMEOUT"`
	ServerWriteTimeout time.Duration `mapstructure:"SERVER_WRITETIMEOUT"`
	ServerIdleTimeout  time.Duration `mapstructure:"SERVER_IDLETIMEOUT"`

	// Database config
	DatabaseURI             string        `mapstructure:"DATABASE_URI"`
	DatabaseMaxConnections  int           `mapstructure:"DATABASE_MAXCONNECTIONS"`
	DatabaseMinConnections  int           `mapstructure:"DATABASE_MINCONNECTIONS"`
	DatabaseMaxConnLifetime time.Duration `mapstructure:"DATABASE_MAXCONNLIFETIME"`
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
	_ = v.ReadInConfig()

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
	v.SetDefault("SERVER_ENV", "development")
	v.SetDefault("SERVER_READTIMEOUT", 45*time.Second)
	v.SetDefault("SERVER_WRITETIMEOUT", 30*time.Second)
	v.SetDefault("SERVER_IDLETIMEOUT", time.Minute)

	// Database defaults
	v.SetDefault("DATABASE_MAXCONNECTIONS", 25)
	v.SetDefault("DATABASE_MINCONNECTIONS", 5)
	v.SetDefault("DATABASE_MAXCONNLIFETIME", 30*time.Minute)
}

func (c *Config) Validate() error {
	// Required fields
	if c.DatabaseURI == "" {
		return fmt.Errorf("DATABASE_URI is required but not set")
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
