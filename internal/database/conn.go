package database

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mbeka02/ticketing-service/config"
)

type Store struct {
	db *pgxpool.Pool
	*Queries
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

		storeInstance = &Store{db: pool, Queries: New(pool)}
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

// ExecTx executes queries within a db transaction
func (s *Store) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	q := s.Queries.WithTx(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *Store) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	err := s.db.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"

	poolStats := s.db.Stat()

	stats["acquire_count"] = strconv.FormatInt(poolStats.AcquireCount(), 10)
	stats["acquired_conns"] = strconv.Itoa(int(poolStats.AcquiredConns()))
	stats["canceled_acquire_count"] = strconv.FormatInt(poolStats.CanceledAcquireCount(), 10)
	stats["constructing_conns"] = strconv.Itoa(int(poolStats.ConstructingConns()))
	stats["empty_acquire_count"] = strconv.FormatInt(poolStats.EmptyAcquireCount(), 10)
	stats["idle_conns"] = strconv.Itoa(int(poolStats.IdleConns()))
	stats["max_conns"] = strconv.Itoa(int(poolStats.MaxConns()))
	stats["total_conns"] = strconv.Itoa(int(poolStats.TotalConns()))

	if poolStats.AcquiredConns() > int32(float64(poolStats.MaxConns())*0.8) {
		stats["message"] = "The database is experiencing heavy load."
	}

	if poolStats.EmptyAcquireCount() > 1000 {
		stats["message"] = "The database has a high number of empty acquire events, indicating potential bottlenecks."
	}

	if poolStats.CanceledAcquireCount() > 100 {
		stats["message"] = "Many connection acquisitions are being canceled, consider increasing pool size or timeout."
	}

	if poolStats.IdleConns() == 0 && poolStats.AcquiredConns() == poolStats.MaxConns() {
		stats["message"] = "Connection pool is exhausted, consider increasing max connections."
	}

	return stats
}
