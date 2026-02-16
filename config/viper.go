package config

import (
	"fmt"
	"strings"
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

// DatabaseConfig is a helper type for passing database configuration
type DatabaseConfig struct {
	URI             string
	MaxConnections  int
	MinConnections  int
	MaxConnLifetime time.Duration
}

// GetDatabaseConfig extracts database config from the main config
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

	v.SetConfigFile(".env")
	v.SetConfigType("env")

	if err := v.ReadInConfig(); err != nil {
		// Only error if file exists but can't be read
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Bind environment variables
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

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
	// Validate required fields
	if c.DatabaseURI == "" {
		return fmt.Errorf("database.uri (DB_URI) is required but not set")
	}

	// Validate constraints
	if c.DatabaseMaxConnections < c.DatabaseMinConnections {
		return fmt.Errorf("database.maxconnections must be >= database.minconnections")
	}

	if c.DatabaseMaxConnections < 1 {
		return fmt.Errorf("database.maxconnections must be at least 1")
	}

	return nil
}
