package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mbeka02/ticketing-service/config"
)

type Store struct {
	db *pgxpool.Pool
}

var (
	storeInstance *Store
	storeOnce     sync.Once
	storeErr      error
)

func NewStore(ctx context.Context, cfg config.DatabaseConfig) (*Store, error) {
	storeOnce.Do(func() {
		poolConfig, err := pgxpool.ParseConfig(cfg.URI)
		if err != nil {
			storeErr = fmt.Errorf("invalid database URI: %w", err)
			return
		}

		poolConfig.MaxConns = int32(cfg.MaxConnections)
		poolConfig.MinConns = int32(cfg.MinConnections)
		poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
		poolConfig.MaxConnIdleTime = 5 * time.Minute
		poolConfig.HealthCheckPeriod = time.Minute

		pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
		if err != nil {
			storeErr = fmt.Errorf("unable to create connection pool: %w", err)
			return
		}

		storeInstance = &Store{db: pool}
	})

	return storeInstance, storeErr
}

func (s *Store) Ping(ctx context.Context) error {
	return s.db.Ping(ctx)
}

func (s *Store) Close() {
	if s.db != nil {
		s.db.Close()
	}
}
